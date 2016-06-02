package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/redhat-developer/henge/pkg/generate/dockercompose"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/runtime"
)

func convertToVersion(objs []runtime.Object, version string) ([]runtime.Object, error) {
	ret := []runtime.Object{}

	for _, obj := range objs {

		convertedObject, err := kapi.Scheme.ConvertToVersion(obj, version)
		if err != nil {
			return nil, err
		}

		ret = append(ret, convertedObject)
	}

	return ret, nil
}

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

	flag.Parse()

	files := flag.Args()
	err := ifFileExists(files)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	template, err := dockercompose.Generate(files...)
	if err != nil {
		return
	}

	var convErr error
	template.Objects, convErr = convertToVersion(template.Objects, "v1")
	if convErr != nil {
		panic(convErr)
	}

	// make it List instead of Template
	list := &kapi.List{Items: template.Objects}

	p, _, err := kubectl.GetPrinter("yaml", "")
	if err != nil {
		panic(err)
	}
	version := unversioned.GroupVersion{Group: "", Version: "v1"}
	p = kubectl.NewVersionedPrinter(p, kapi.Scheme, version)
	p.PrintObj(list, os.Stdout)

}
