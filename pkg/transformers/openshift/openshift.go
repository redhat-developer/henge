package openshift

import (
	"os"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/openshift/origin/pkg/generate/dockercompose"
	templateapi "github.com/openshift/origin/pkg/template/api"

	// Install OpenShift APIs
	_ "github.com/openshift/origin/pkg/build/api/install"
	_ "github.com/openshift/origin/pkg/deploy/api/install"
	_ "github.com/openshift/origin/pkg/image/api/install"
	_ "github.com/openshift/origin/pkg/route/api/install"
	_ "github.com/openshift/origin/pkg/template/api/install"
)

func Transform(paths ...string) (*templateapi.Template, error) {

	template, err := dockercompose.Generate(paths...)

	if err != nil {
		return nil, err
	}

	// Convert template objects to versioned objects
	var convErr error
	template.Objects, convErr = convertToVersion(template.Objects, "v1")
	if convErr != nil {
		panic(convErr)
	}

	return template, err
}

// Print openshift template
func Print(template *templateapi.Template) {
	// make it List instead of Template
	list := &kapi.List{Items: template.Objects}

	printer, _, _err := kubectl.GetPrinter("yaml", "")
	if _err != nil {
		panic(_err)
	}
	version := unversioned.GroupVersion{Group: "", Version: "v1"}
	printer = kubectl.NewVersionedPrinter(printer, kapi.Scheme, version)
	printer.PrintObj(list, os.Stdout)
}

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
