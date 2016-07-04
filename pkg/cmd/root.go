package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/redhat-developer/henge/pkg/types"
)

const cliLong = `
Henge takes a docker-compose.yml file and converts it into a openshift or kubernetes artifacts,
which then can be used to deploy apps on that platforms.
`

const example = `
To convert the docker-compose.yml file in the current directory to openshift's artifacts

  $ henge openshift

To convert the file of your choice to kubernetes's artifacts.

  $ henge kubernetes -f foo.yml

To convert docker-compose.yml file in current directory and also ask questions interactively.

  $ henge openshift -i

To provide multiple file for conversion

  $ henge kubernetes -f foo.yml,bar.yml,docker-compose.yml
`

func Execute() (*types.CmdValues, error) {
	var val types.CmdValues

	var RootCmd = &cobra.Command{
		Use:     "henge",
		Short:   "Henge converts the docker compose file to various orchestration providers' artifacts.",
		Long:    cliLong,
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Target not specified")
			}
			val.Target = args[0]
			return nil
		},
	}
	RootCmd.Flags().BoolVarP(&val.Interactive, "interactive", "i", false, "Ask questions about values that can affect conversion.")
	RootCmd.Flags().StringSliceVarP(&val.Files, "files", "f", []string{"docker-compose.yml"}, "Provide docker-compose files, comma separated.")
	RootCmd.Flags().IntVarP(&val.Loglevel, "loglevel", "", 0, "Log level to show.")

	if err := RootCmd.Execute(); err != nil {
		return nil, err
	}
	return &val, nil
}
