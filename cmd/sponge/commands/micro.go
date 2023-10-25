package commands

import (
	"github.com/i2dou/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// GenMicroCommand generate micro service code
func GenMicroCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "micro",
		Short:         "Generate proto, model, cache, dao, service, rpc, rpc-gw, rpc-cli code",
		Long:          "generate proto, model, cache, dao, service, rpc, rpc-gw, rpc-cli code.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		generate.ProtoBufCommand(),
		generate.ModelCommand("micro"),
		generate.DaoCommand("micro"),
		generate.CacheCommand("micro"),
		generate.ServiceCommand(),
		generate.RPCCommand(),
		generate.RPCGwPbCommand(),
		generate.RPCPbCommand(),
		generate.RPCConnectionCommand(),
		generate.ConvertSwagJSONCommand("micro"),
	)

	return cmd
}
