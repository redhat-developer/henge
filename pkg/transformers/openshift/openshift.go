package openshift

import (
	kapi "k8s.io/kubernetes/pkg/api"

	templateapi "github.com/openshift/origin/pkg/template/api"

	// Install OpenShift APIs
	_ "github.com/openshift/origin/pkg/build/api/install"
	_ "github.com/openshift/origin/pkg/deploy/api/install"
	_ "github.com/openshift/origin/pkg/image/api/install"
	_ "github.com/openshift/origin/pkg/route/api/install"
	_ "github.com/openshift/origin/pkg/template/api/install"
)

func Transform(interactive bool, paths ...string) (*kapi.List, error) {

	template, err := Generate(interactive, paths...)

	if err != nil {
		return nil, err
	}

	list := ConvertToList(template)

	return list, err
}

// Convert OpenShift Template to Kubernetes List
func ConvertToList(template *templateapi.Template) *kapi.List {
	list := &kapi.List{Items: template.Objects}
	return list
}
