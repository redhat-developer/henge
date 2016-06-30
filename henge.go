package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/redhat-developer/henge/pkg/transformers"
	"github.com/redhat-developer/henge/pkg/utils"
)

func main() {
	target := flag.String("target", "", "Target platform (openshift or kubernetes)")
	interactive := flag.Bool("interactive", false, "Ask questions about missing arguments.")

	flag.Parse()

	if *target == "" {
		fmt.Fprintln(os.Stderr, "You must provide target platform using -target argument.")
		os.Exit(1)
	}

	files := flag.Args()
	err := utils.CheckIfFileExists(files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = transformers.Transform(*target, *interactive, flag.Args()[0:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
