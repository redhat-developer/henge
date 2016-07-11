package main

import (
	"os"

	"github.com/redhat-developer/henge/pkg/cmd"
)

func main() {
	// parse all command line args
	_, err := cmd.Execute()
	if err != nil {
		//fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
