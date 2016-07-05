package transformers

import (
	"fmt"

	"github.com/redhat-developer/henge/pkg/types"

	"github.com/redhat-developer/henge/pkg/transformers/kubernetes"
	"github.com/redhat-developer/henge/pkg/transformers/openshift"
)

// Transform transforms the given project into artifacts for the specified
// orchestration provider
func Transform(vals *types.CmdValues) error {
	switch vals.Target {
	case "openshift":
		list, err := openshift.Transform(vals)
		if err != nil {
			return err
		}
		kubernetes.PrintList(list)
		return nil
	case "kubernetes":
		list, err := kubernetes.Transform(vals)
		if err != nil {
			return err
		}
		kubernetes.PrintList(list)
		return nil
	}
	err := fmt.Errorf("Provider not supported")
	return err
}
