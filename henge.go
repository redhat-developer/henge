package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rtnpro/henge/pkg/transformers"
)

func main() {
	provider := flag.String("provider", "openshift", "Target provider")

	flag.Parse()

	files := flag.Args()
	err := ifFileExists(files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err := transformers.Transform(*provider, flag.Args()[0:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
