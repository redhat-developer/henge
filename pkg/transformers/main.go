package transformers

import (
	"fmt"

	"github.com/redhat-developer/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// provider
func Transform(provider string, paths ...string) error {
	if provider == "openshift" {
		obj, err := openshift.Transform(paths...)
		if err != nil {
			return err
		}
		openshift.Print(obj)
	} else {
		err := fmt.Errorf("Provider not supported")
		return err
	}
	return nil
}
