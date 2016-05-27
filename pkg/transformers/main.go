package transformers

import (
	"fmt"
	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"

	// "github.com/rtnpro/henge/pkg/transformers/kubernetes"
	// "github.com/rtnpro/henge/pkg/transformers/marathon"
	"github.com/rtnpro/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// provider
func Transform(provider string, project project.Project, bases []string) error {
	if provider == "openshift" {
		err := openshift.Transform(project, bases)
		if err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("Provider not supported")
		return err
	}
	return nil
}
