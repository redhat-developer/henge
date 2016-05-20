package main

import (
	"flag"
	"os"

	"github.com/rtnpro/henge/pkg/generate/dockercompose"

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

func main() {

	flag.Parse()

	template, err := dockercompose.Generate(flag.Args()[0:]...)
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
