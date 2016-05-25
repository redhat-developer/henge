package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rtnpro/henge/pkg/loaders/compose"
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

	project, err = compose.Load(flag.Args()[0:]...)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	transformers.Transform(provider, project)
}
