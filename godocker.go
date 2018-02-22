package godocker

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type Container struct {
	ctx    context.Context // docker context
	cli    *client.Client  // docker client
	id     int             // Container id, once started
	indir  string          // Input directory that will be bound to /indir image dir
	outdir string          // Output directory that will be bound in /outdir image dir
	cl     []string        // Command line that will be run in the container
	name   string          // Image name
}

// Creates a new docker container, pull the docker image if necessary
func NewContainer(name string) (c *Container, err error) {
	var ctx context.Context
	var cli *client.Client
	var reader io.ReadCloser

	ctx = context.Background()
	if cli, err = client.NewEnvClient(); err != nil {
		return
	}

	c = &Container{
		ctx:    ctx,
		cli:    cli,
		id:     -1,
		indir:  "",
		outdir: "",
		cl:     []string{},
		name:   name,
	}

	if reader, err = cli.ImagePull(ctx, name, types.ImagePullOptions{}); err != nil {
		return
	}

	io.Copy(os.Stdout, reader)

	return
}

// Sets the input host directory to be bound to /indir image directory
func (c *Container) SetInputDir(indir string) {
	c.indir = indir
}

// Sets the output host directory to be bound to /outdir image directory
func (c *Container) SetOutputDir(outdir string) {
	c.outdir = outdir
}

func (c *Container) Start() (err error) {
	var binds []string
	var body container.ContainerCreateCreatedBody

	// To run commands as current os user in the container
	// That way output files are owned by current user
	var u *user.User
	if u, err = user.Current(); err != nil {
		return
	}

	if len(c.cl) == 0 {
		return errors.New("Cannot start container, command line is empty")
	}

	if c.indir == "" {
		return errors.New("Cannot start container, do not know which input directory to bind to /indir")
	}
	if _, err := os.Stat(c.indir); os.IsNotExist(err) {
		return errors.New("Cannot start container, given indir does not exist")
	}

	if c.outdir == "" {
		return errors.New("Cannot start container, do not know which output directory to bind to /outdir")
	}
	if _, err := os.Stat(c.outdir); os.IsNotExist(err) {
		return errors.New("Cannot start container, given outdir does not exist")
	}

	binds = make([]string, 0)
	binds = append(binds, fmt.Sprintf("%s:/indir", c.indir))
	binds = append(binds, fmt.Sprintf("%s:/outdir", c.outdir))

	body, err = c.cli.ContainerCreate(
		c.ctx,
		&container.Config{
			Image:      c.name,
			Cmd:        c.cl,
			Tty:        true,
			User:       fmt.Sprintf("%s:%s", u.Uid, u.Gid),
			WorkingDir: "/indir",
			Entrypoint: make([]string, 0)},
		&container.HostConfig{
			Binds: binds,
		},
		nil,
		"")

	if err != nil {
		return
	}

	if err = c.cli.ContainerStart(c.ctx, body.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	// We intercept signals to stop the container at the same time
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		d := 5 * time.Second
		for sig := range sigchan {
			log.Print(sig)
			if err := c.cli.ContainerStop(c.ctx, body.ID, &d); err != nil {
				log.Print(err)
			}
			os.Exit(1)
		}
	}()

	// We wait for the end of the execution of the command
	statusCh, errCh := c.cli.ContainerWait(c.ctx, body.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	// We get the logs of the execution on stdout
	out, err := c.cli.ContainerLogs(c.ctx, body.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	return
}

func (c *Container) SetCommandLine(cl []string) {
	c.cl = append(c.cl, cl...)
}
