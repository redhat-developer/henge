package kubernetes

import (
	"os"

	"github.com/redhat-developer/henge/pkg/types"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/runtime"
)

func Transform(vals *types.CmdValues) (*kapi.List, error) {
	list, err := Generate(vals)
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

// PrintList will either print List in Yaml format to stdout or
// a file, on the location given on commandline
func PrintList(list *kapi.List, vals *types.CmdValues) error {
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

	out := os.Stdout
	if vals.OutputFile != "" {
		f, err := os.Create(vals.OutputFile)
		if err != nil {
			return err
		}
		out = f
		defer f.Close()
	}
	if err = p.PrintObj(list, out); err != nil {
		return err
	}
	return nil
}
