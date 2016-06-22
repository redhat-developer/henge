package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/redhat-developer/henge/pkg/transformers"
)

// Loop over a array of filepaths and check if it exists
// if it exists check if it is not a directory.
func ifFileExists(files []string) error {
	for _, filename := range files {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("main: file %q not found", filename)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("main: %q is a directory", filename)
		}
	}
	return nil
}

func main() {
	provider := flag.String("target", "", "Target platform (openshift or kubernetes)")

	flag.Parse()

	if *provider == "" {
		fmt.Fprintln(os.Stderr, "You must provide target platform using -platform argument.")
		os.Exit(1)
	}

	files := flag.Args()
	err := ifFileExists(files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = transformers.Transform(*provider, flag.Args()[0:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
