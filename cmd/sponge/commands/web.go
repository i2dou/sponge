package commands

import (
	"github.com/i2dou/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// GenWebCommand generate web server code
func GenWebCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "web",
		Short:         "Generate model, cache, dao, handler, web code",
		Long:          "generate model, cache, dao, handler, web code.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		generate.ModelCommand("web"),
		generate.DaoCommand("web"),
		generate.CacheCommand("web"),
		generate.HandlerCommand(),
		generate.HTTPCommand(),
		generate.HTTPPbCommand(),
		generate.ConvertSwagJSONCommand("web"),
		generate.HandlerPbCommand(),
	)

	return cmd
}
