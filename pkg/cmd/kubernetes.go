package cmd

import (
	"fmt"
	"github.com/redhat-developer/henge/pkg/transformers/kubernetes"
	"github.com/redhat-developer/henge/pkg/types"
	"github.com/spf13/cobra"
	"os"
)

func kubernetesCmd(vals *types.CmdValues) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "convert to Kubernetes artifacts",
		Long:  "To convert the docker-compose.yml file in the current directory to kubernetes' artifacts",
		Run: func(cmd *cobra.Command, args []string) {

			errorIfFileDoesNotExist(vals)

			list, err := kubernetes.Transform(vals)

			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			err = kubernetes.PrintList(list, vals)

			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
	}

	addProviderFlags(cmd, vals)
	return cmd
}
