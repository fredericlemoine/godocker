package godocker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	RAXML_GTRGAMMA = iota
	RAXML_GTRCAT
	RAXML_DAYHOFF
	RAXML_DCMUT
	RAXML_JTT
	RAXML_MTREV
	RAXML_WAG
	RAXML_RTREV
	RAXML_CPREV
	RAXML_VT
	RAXML_BLOSUM62
	RAXML_MTMAM
	RAXML_LG
	RAXML_MTART
	RAXML_MTZOA
	RAXML_PMB
	RAXML_HIVB
	RAXML_HIVW
	RAXML_JTTDCMUT
	RAXML_FLU
	RAXML_STMTREV
	RAXML_DUMMY
	RAXML_DUMMY2
	RAXML_LG4M
	RAXML_LG4X
	RAXML_GTR_UNLINKED
	RAXML_GTR
)

const (
	RAXML_IMAGE_V8_2_11 = iota
)

var RAxMLModels = []string{
	"GTRGAMMA",
	"GTRCAT",
	"PROTGAMMADAYHOFF",
	"PROTGAMMADCMUT",
	"PROTGAMMAJTT",
	"PROTGAMMAMTREV",
	"PROTGAMMAWAG",
	"PROTGAMMARTREV",
	"PROTGAMMACPREV",
	"PROTGAMMAVT",
	"PROTGAMMABLOSUM62",
	"PROTGAMMAMTMAM",
	"PROTGAMMALG",
	"PROTGAMMAMTART",
	"PROTGAMMAMTZOA",
	"PROTGAMMAPMB",
	"PROTGAMMAHIVB",
	"PROTGAMMAHIVW",
	"PROTGAMMAJTTDCMUT",
	"PROTGAMMAFLU",
	"PROTGAMMASTMTREV",
	"PROTGAMMADUMMY",
	"PROTGAMMADUMMY2",
	"PROTGAMMALG4M",
	"PROTGAMMALG4X",
	"PROTGAMMAGTR_UNLINKED",
	"PROTGAMMAGTR",
}

var RAxMLImages = []string{
	"docker.io/evolbioinfo/raxml:v8.2.11",
}

type RAxMLTool struct {
	inputalign string
	outputtree string
	model      int
	gammacat   int
	cpus       int
	seed       int
	image      int // RAXML_IMAGE
	runname    string
	force      bool // if true, will remove previous run files if exists
}

func NewRAXMLTool() *RAxMLTool {
	return &RAxMLTool{
		model:    RAXML_GTRGAMMA,
		image:    RAXML_IMAGE_V8_2_11,
		seed:     12345,
		gammacat: 4,
		runname:  "DOCKER",
		force:    false,
	}
}

func (r *RAxMLTool) SetRunName(name string) {
	r.runname = name

}

func (r *RAxMLTool) RunName() string {
	return r.runname
}

func (r *RAxMLTool) SetForce(force bool) {
	r.force = force
}

func (r *RAxMLTool) Force() bool {
	return r.force
}

// if cpus <0, then it will be automatically computed
// be the function Cpus()
func (r *RAxMLTool) SetCpus(cpus int) {
	r.cpus = cpus
}

// Compute the  number of cpus that can be used by
// the process.
//
// if r.cpus <0 or > max number of cpus on the machine,
// then returns maxnumber of cpus -1
func (r *RAxMLTool) Cpus() (cpus int) {
	if r.cpus <= 0 || r.cpus > runtime.NumCPU() {
		cpus = runtime.NumCPU() - 1
		if cpus <= 0 {
			cpus = 1
		}
	} else {
		cpus = r.cpus
	}
	return
}

func (r *RAxMLTool) SetInputAlign(align string) (err error) {
	r.inputalign = align
	_, err = os.Stat(align)
	return
}

func (r *RAxMLTool) InputAlign() string {
	return r.inputalign
}

func (r *RAxMLTool) SetOutputTree(tree string) (err error) {
	var outdir string
	r.outputtree = tree
	outdir, err = filepath.Abs(filepath.Dir(r.outputtree))
	if err == nil {
		_, err = os.Stat(outdir)
	}
	return
}

func (r *RAxMLTool) OutputTree() string {
	return r.outputtree
}

func (r *RAxMLTool) InDir() (indir string, err error) {
	indir, err = filepath.Abs(filepath.Dir(r.inputalign))
	if err == nil {
		_, err = os.Stat(indir)
	}
	return
}

func (r *RAxMLTool) InBaseName() (basename string) {
	return filepath.Base(r.inputalign)
}

func (r *RAxMLTool) OutDir() (outdir string, err error) {
	outdir, err = filepath.Abs(filepath.Dir(r.outputtree))
	if err == nil {
		_, err = os.Stat(outdir)
	}
	return
}

func (r *RAxMLTool) ModelString() (model string, err error) {
	if r.model < len(RAxMLModels) && r.model >= 0 {
		model = RAxMLModels[r.model]
	} else {
		err = errors.New(fmt.Sprintf("Invalid Model: %d", r.model))
	}
	return
}

func (r *RAxMLTool) SetModelString(model string) (err error) {
	for i, m := range RAxMLModels {
		if m == model {
			r.model = i
			return
		}
	}
	err = errors.New(fmt.Sprintf("No such model: %s", model))
	return
}
func (r *RAxMLTool) SetModel(model int) (err error) {
	if model < 0 || model > len(RAxMLModels) {
		err = errors.New(fmt.Sprintf("No such model: %d", model))
	}

	r.model = model
	return
}

func (r *RAxMLTool) SetImage(image int) (err error) {
	if image < 0 || image > len(RAxMLImages) {
		err = errors.New(fmt.Sprintf("No such image: %d", image))
	}

	r.image = image
	return
}

func (r *RAxMLTool) SetImageString(image string) (err error) {
	for i, v := range RAxMLImages {
		if v == image {
			r.image = i
			return
		}
	}
	err = errors.New(fmt.Sprintf("No such image: %s", image))
	return
}

func (r *RAxMLTool) ImageString() (image string, err error) {
	if r.image < len(RAxMLImages) && r.image >= 0 {
		image = RAxMLImages[r.image]
	} else {
		err = errors.New(fmt.Sprintf("Invalid image:  %d", r.image))
	}
	return
}

// Executes the current RAxML command
//
// Will start a docker container, set the input and output
// will run the command, and move the resulting tree in
// the given file
func (r *RAxMLTool) Execute() (err error) {
	var c *Container
	var indir, outdir string
	var image string
	var raxmloutput string
	var cl []string

	if image, err = r.ImageString(); err != nil {
		return
	}
	if c, err = NewContainer(image); err != nil {
		return
	}
	if indir, err = r.InDir(); err != nil {
		return
	}
	if outdir, err = r.OutDir(); err != nil {
		return
	}
	c.SetInputDir(indir)
	c.SetOutputDir(outdir)

	if cl, err = r.CommandLine(); err != nil {
		return
	}
	c.SetCommandLine(cl)
	if r.Force() {
		if err = r.RemoveIfExistsPrevRun(); err != nil {
			return
		}
	}
	if err = c.Start(); err != nil {
		return
	}

	raxmloutput = fmt.Sprintf("%s/RAxML_bestTree.TEST", indir)

	// When execution is over, we move the output file to the given name
	os.Rename(raxmloutput, r.OutputTree())

	return
}

// Builds the command line that will run RAxML
func (r *RAxMLTool) CommandLine() (cl []string, err error) {
	var model string

	if model, err = r.ModelString(); err != nil {
		return
	}

	cl = []string{"raxmlHPC",
		"-f", "d",
		"-m", model,
		"-c", fmt.Sprintf("%d", r.gammacat),
		"-s", r.InBaseName(),
		"-n", r.RunName(),
		"-T", fmt.Sprintf("%d", r.Cpus()),
		"-p", fmt.Sprintf("%d", r.seed),
	}
	return
}

// Will remove run files from a previous run with the same name
// if it exists
func (r *RAxMLTool) RemoveIfExistsPrevRun() (err error) {
	var indir string

	if indir, err = r.InDir(); err != nil {
		return
	}

	if _, err2 := os.Stat(fmt.Sprintf("%s/RAxML_info.%s", indir, r.RunName())); err2 == nil {
		os.Remove(fmt.Sprintf("%s/RAxML_bestTree.%s", indir, r.RunName()))
		os.Remove(fmt.Sprintf("%s/RAxML_info.%s", indir, r.RunName()))
		os.Remove(fmt.Sprintf("%s/RAxML_log.%s", indir, r.RunName()))
		os.Remove(fmt.Sprintf("%s/RAxML_parsimonyTree.%s", indir, r.RunName()))
		os.Remove(fmt.Sprintf("%s/RAxML_result.%s", indir, r.RunName()))
	}
	return
}
