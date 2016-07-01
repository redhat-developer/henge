package main

import (
	"fmt"
	"os"

	"github.com/redhat-developer/henge/pkg/cmd"
	"github.com/redhat-developer/henge/pkg/transformers"
	"github.com/redhat-developer/henge/pkg/utils"
)

func main() {
	// parse all command line args
	vals, err := cmd.Execute()
	if err != nil {
		//fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// check if files exists
	err = utils.CheckIfFileExists(vals.Files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// make conversion
	err = transformers.Transform(vals.Target, vals.Interactive, vals.Files...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
