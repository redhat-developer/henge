package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/redhat-developer/henge/pkg/types"
	"github.com/redhat-developer/henge/pkg/utils"
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
	}

	RootCmd.PersistentFlags().BoolVarP(&val.Interactive, "interactive", "i", false, "Ask questions about values that can affect conversion.")
	RootCmd.PersistentFlags().IntVarP(&val.Loglevel, "loglevel", "", 0, "Log level to show.")

	RootCmd.AddCommand(openshiftCmd(&val))
	RootCmd.AddCommand(kubernetesCmd(&val))

	if err := RootCmd.Execute(); err != nil {
		return nil, err
	}

	return &val, nil
}

func addProviderFlags(cmd *cobra.Command, vals *types.CmdValues) {

	cmd.Flags().StringSliceVarP(&vals.Files, "files", "f", nil, "Provide comma separated docker-compose files. Default: 'docker-compose.yml' or 'docker-compose.yaml'")
	cmd.Flags().StringVarP(&vals.OutputFile, "output-file", "o", "", "File to save converted artifacts.")
}

func fileDefaultsAndSanity(val *types.CmdValues) error {
	// check if files exists

	if val.Files == nil {
		val.Files = []string{"docker-compose.yml"}
		err := utils.CheckIfFileExists(val.Files)

		if err == nil {
			return nil
		}

		val.Files = []string{"docker-compose.yaml"}
		err = utils.CheckIfFileExists(val.Files)

		if err == nil {
			return nil
		}

		return fmt.Errorf("docker-compose.yml or docker-compose.yaml file not found\n")
	}

	return utils.CheckIfFileExists(val.Files)
}
