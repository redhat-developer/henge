package kubernetes

import (
	"os"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/runtime"
)

func Transform(interactive bool, paths ...string) (*kapi.List, error) {
	list, err := Generate(interactive, paths...)
	return list, err
}

// Convert all objects in objs to versioned objects
func ConvertToVersion(objs []runtime.Object, version string) ([]runtime.Object, error) {
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

//Display List in Yaml format on stdout
func PrintList(list *kapi.List) {
	var convErr error
	version := unversioned.GroupVersion{Group: "", Version: "v1"}

	list.Items, convErr = ConvertToVersion(list.Items, version.Version)
	if convErr != nil {
		panic(convErr)
	}

	p, _, err := kubectl.GetPrinter("yaml", "")
	if err != nil {
		panic(err)
	}
	p = kubectl.NewVersionedPrinter(p, kapi.Scheme, version)
	p.PrintObj(list, os.Stdout)
}
