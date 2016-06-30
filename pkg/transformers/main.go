package transformers

import (
	"fmt"

	"github.com/redhat-developer/henge/pkg/transformers/kubernetes"
	"github.com/redhat-developer/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// orchestration provider
func Transform(target string, interactive bool, paths ...string) error {
	switch target {
	case "openshift":
		list, err := openshift.Transform(interactive, paths...)
		if err != nil {
			return err
		}
		kubernetes.PrintList(list)
		return nil
	case "kubernetes":
		list, err := kubernetes.Transform(interactive, paths...)
		if err != nil {
			return err
		}
		kubernetes.PrintList(list)
		return nil
	}
	err := fmt.Errorf("Provider not supported")
	return err
}
