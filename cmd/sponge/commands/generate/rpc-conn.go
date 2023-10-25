package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/i2dou/sponge/pkg/replacer"

	"github.com/spf13/cobra"
)

// RPCConnectionCommand generate rpc connection code
func RPCConnectionCommand() *cobra.Command {
	var (
		moduleName     string // module name for go.mod
		outPath        string // output directory
		rpcServerNames string // rpc service names
	)

	cmd := &cobra.Command{
		Use:   "rpc-conn",
		Short: "Generate rpc connection code",
		Long: `generate rpc connection code.

Examples:
  # generate rpc connection code
  sponge micro rpc-conn --module-name=yourModuleName --rpc-server-name=user

  # generate rpc connection code with multiple names.
  sponge micro rpc-conn --module-name=yourModuleName --rpc-server-name=name1,name2

  # generate rpc connection code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge micro rpc-conn --rpc-server-name=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, _ := getNamesFromOutDir(outPath)
			if mdName != "" {
				moduleName = mdName
			} else if moduleName == "" {
				return errors.New(`required flag(s) "module-name" not set, use "sponge micro rpc-conn -h" for help`)
			}

			rpcNames := strings.Split(rpcServerNames, ",")
			for _, rpcName := range rpcNames {
				if rpcName == "" {
					continue
				}

				var err error
				outPath, err = runGenRPCConnectionCommand(moduleName, rpcName, outPath)
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  move the folder "internal" to your project code folder.

`)
			fmt.Printf("generate \"rpc-conn\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	cmd.Flags().StringVarP(&rpcServerNames, "rpc-server-name", "r", "", "rpc service name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("rpc-server-name")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./rpc-conn_<time>,"+
		" if you specify the directory where the web or microservice generated by sponge, the module-name flag can be ignored")

	return cmd
}

func runGenRPCConnectionCommand(moduleName string, rpcName string, outPath string) (string, error) {
	subTplName := "rpc-conn"
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{ // only the specified subdirectory is processed, if empty or no subdirectory is specified, it means all files
		"internal/rpcclient",
	}
	ignoreDirs := []string{} // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{ // specify the files in the subdirectory to be ignored for processing
		"doc.go", "serverNameExample_test.go",
	}

	r.SetSubDirsAndFiles(subDirs)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addRPCConnectionFields(moduleName, rpcName)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func addRPCConnectionFields(moduleName string, serverName string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, []replacer.Field{
		{
			Old: "github.com/i2dou/sponge/configs",
			New: moduleName + "/configs",
		},
		{
			Old: "github.com/i2dou/sponge/internal/config",
			New: moduleName + "/internal/config",
		},
		{
			Old:             "serverNameExample",
			New:             serverName,
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
