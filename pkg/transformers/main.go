package transformers

import (
	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"

	"github.com/rtnpro/henge/pkg/transformers/kubernetes"
	"github.com/rtnpro/henge/pkg/transformers/marathon"
	"github.com/rtnpro/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// provider
func Transform(provider string, project *project.Project) error {
	if provider == "kubernetes" {
		kubernetes.Transform(project)
	} else if provider == "openshift" {
		openshift.Transform(project)
	} else if provider == "marathon" {
		marathon.Transform(project)
	} else {
		err := "Provider not supported"
		return err
	}
}
