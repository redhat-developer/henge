package main

import (
	"flag"
	"fmt"

	"github.com/kadel/henge/pkg/generate/dockercompose"
)

func main() {

	flag.Parse()

	template, err := dockercompose.Generate(flag.Args()[0:]...)
	if err != nil {
		return
	}

	fmt.Printf("%#v", template)

}
