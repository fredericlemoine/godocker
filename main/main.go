package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fredericlemoine/godocker"
)

func main() {
	var input, output string
	var force bool

	flag.StringVar(&input, "infile", "none", "Input alignment file")
	flag.StringVar(&output, "outfile", "none", "Output alignment file")
	flag.BoolVar(&force, "force", false, "Remove previous runs if they exist")

	flag.Parse()

	log.Print(fmt.Sprintf("Input file : %s", input))
	log.Print(fmt.Sprintf("Output file : %s", output))

	if input == "none" {
		log.Fatal("input file must be provided")
	}
	if output == "none" {
		log.Fatal("output file must be provided")
	}
	rx := godocker.NewRAXMLTool()
	rx.SetForce(force)
	rx.SetCpus(-1)
	rx.SetInputAlign(input)
	rx.SetOutputTree(output)
	log.Print(rx.Execute())
}
