package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/i2dou/sponge/pkg/replacer"
	"github.com/i2dou/sponge/pkg/sql2code"
	"github.com/i2dou/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// ModelCommand generate model code
func ModelCommand(parentName string) *cobra.Command {
	var (
		outPath  string // output directory
		dbTables string // table names

		sqlArgs = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate model code based on mysql table",
		Long: fmt.Sprintf(`generate model code based on mysql table.

Examples:
  # generate model code and embed gorm.model struct.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate model code with multiple table names.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # generate model code, structure fields correspond to the column names of the table.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate model code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`, parentName, parentName, parentName, parentName),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			tableNames := strings.Split(dbTables, ",")
			for _, tableName := range tableNames {
				if tableName == "" {
					continue
				}

				sqlArgs.DBTable = tableName
				codes, err := sql2code.Generate(&sqlArgs)
				if err != nil {
					return err
				}

				outPath, err = runGenModelCommand(codes, outPath)
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  move the folder "internal" to your project code folder.

`)
			fmt.Printf("generate \"model\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed gorm.model struct")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-name-type", "j", 1, "json tags name type, 0:snake case, 1:camel case")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./model_<time>")

	return cmd
}

func runGenModelCommand(codes map[string]string, outPath string) (string, error) {
	subTplName := "model"
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{"internal/model"} // only the specified subdirectory is processed, if empty or no subdirectory is specified, it means all files
	ignoreDirs := []string{}              // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{              // specify the files in the subdirectory to be ignored for processing
		"init.go", "init_test.go",
	}

	r.SetSubDirsAndFiles(subDirs)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addModelFields(r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func addModelFields(r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
