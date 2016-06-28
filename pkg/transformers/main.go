package transformers

import (
	"fmt"

	"github.com/redhat-developer/henge/pkg/transformers/kubernetes"
	"github.com/redhat-developer/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// provider
func Transform(provider string, interactive bool, paths ...string) error {
	switch provider {
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
