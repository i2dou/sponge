package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// InitCommand initial sponge
func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize sponge",
		Long: `initialize sponge.

Examples:
  # run init, download code and install tools.
  sponge init
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("initialize sponge ......")

			// download sponge template code
			err := runUpgradeCommand()
			if err != nil {
				return err
			}
			ver, err := copyToTempDir()
			if err != nil {
				return err
			}
			updateSpongeInternalPlugin(ver)

			// installing dependent plug-ins
			_, lackNames := checkInstallTools()
			installTools(lackNames)

			return nil
		},
	}

	return cmd
}
