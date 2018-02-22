package godocker

type DockerTool interface {
	Execute() error
	SetCommandLine()
	SetContainer()
}
