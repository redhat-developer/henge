package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rtnpro/henge/pkg/loaders"
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

	project, bases, err := loaders.Compose(flag.Args()[0:]...)

	fmt.Println("provider: ", *provider)
	fmt.Println("project: ", *project)
	fmt.Println("bases: ", bases)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	transformers.Transform(*provider, project, bases)
}
