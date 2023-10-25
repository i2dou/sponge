package generate

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/i2dou/sponge/pkg/gofile"
	"github.com/i2dou/sponge/pkg/replacer"
	"github.com/i2dou/sponge/pkg/sql2code"
	"github.com/i2dou/sponge/pkg/sql2code/parser"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// RPCCommand generate rpc service code
func RPCCommand() *cobra.Command {
	var (
		moduleName  string // module name for go.mod
		serverName  string // server name
		projectName string // project name for deployment name
		repoAddr    string // image repo address
		outPath     string // output directory
		dbTables    string // table names
		sqlArgs     = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "rpc",
		Short: "Generate rpc service code based on mysql table",
		Long: `generate rpc service code based on mysql table.

Examples:
  # generate rpc service code and embed gorm.model struct.
  sponge micro rpc --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate rpc service code, structure fields correspond to the column names of the table.
  sponge micro rpc --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate rpc service code with multiple table names.
  sponge micro rpc --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # generate rpc service code and specify the output directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge micro rpc --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir

  # generate rpc service code and specify the docker image repository address.
  sponge micro rpc --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var firstTable string
			var servicesTableNames []string
			tableNames := strings.Split(dbTables, ",")
			if len(tableNames) == 1 {
				firstTable = tableNames[0]
			} else if len(tableNames) > 1 {
				firstTable = tableNames[0]
				servicesTableNames = tableNames[1:]
			}

			projectName, serverName = convertProjectAndServerName(projectName, serverName)

			sqlArgs.DBTable = firstTable
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}
			outPath, err = runGenRPCCommand(moduleName, serverName, projectName, repoAddr, sqlArgs.DBDsn, codes, outPath)
			if err != nil {
				return err
			}

			for _, serviceTableName := range servicesTableNames {
				if serviceTableName == "" {
					continue
				}

				sqlArgs.DBTable = serviceTableName
				codes, err := sql2code.Generate(&sqlArgs)
				if err != nil {
					return err
				}

				outPath, err = runGenServiceCommand(moduleName, serverName, codes, outPath)
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  1. open a terminal and execute the command to generate code:  make proto
  2. compile and run service:   make run
  3. open the file internal/service/xxx_client_test.go using Goland or VS Code, and test CRUD api interface.

`)
			fmt.Printf("generate %s's rpc service code successfully, out = %s\n", serverName, outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed gorm.model struct")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-name-type", "j", 1, "json tags name type, 0:snake case, 1:camel case")
	cmd.Flags().StringVarP(&repoAddr, "repo-addr", "r", "", "docker image repository address, excluding http and repository names")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./serverName_rpc_<time>")

	return cmd
}

func runGenRPCCommand(moduleName string, serverName string, projectName string, repoAddr string,
	dbDSN string, codes map[string]string, outPath string) (string, error) {
	subTplName := "rpc"
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{ // specify the subdirectory for processing
		"sponge/api", "cmd/serverNameExample_grpcExample", "sponge/configs", "sponge/deployments",
		"sponge/scripts", "sponge/internal", "sponge/third_party",
	}
	subFiles := []string{ // specify the sub-documents to be processed
		"sponge/.gitignore", "sponge/.golangci.yml", "sponge/go.mod", "sponge/go.sum",
		"sponge/Jenkinsfile", "sponge/Makefile", "sponge/README.md",
	}
	ignoreDirs := []string{ // specify the directory in the subdirectory where processing is ignored
		"internal/handler", "internal/rpcclient", "internal/routers", "internal/types",
	}
	ignoreFiles := []string{ // specify the files in the subdirectory to be ignored for processing
		"types.pb.validate.go", "types.pb.go", // api/types
		"userExample.pb.go", "userExample.pb.validate.go", "userExample_grpc.pb.go", "userExample_router.pb.go", // api/serverNameExample/v1
		"userExample_http.go", "systemCode_http.go", // internal/ecode
		"http.go", "http_option.go", "http_test.go", // internal/server
		"userExample_logic.go", "userExample_logic_test.go", "service/userExample_test.go", // internal/service
		"scripts/swag-docs.sh",                                      // sponge/scripts
		"doc.go", "cacheNameExample.go", "cacheNameExample_test.go", // internal/cache
	}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addRPCFields(moduleName, serverName, projectName, repoAddr, r, dbDSN, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}
	_ = saveGenInfo(moduleName, serverName, r.GetOutputDir())

	return r.GetOutputDir(), nil
}

func addRPCFields(moduleName string, serverName string, projectName string, repoAddr string,
	r replacer.Replacer, dbDSN string, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	repoHost, _ := parseImageRepoAddr(repoAddr)

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, protoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, serviceClientFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, serviceTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFileBuild, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerComposeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sDeploymentFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sServiceFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, imageBuildFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, imageBuildLocalFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteAllFieldsMark(r, makeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, gitIgnoreFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteAllFieldsMark(r, protoShellFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, appConfigFile, wellStartMark, wellEndMark)...)
	fields = append(fields, replaceFileContentMark(r, readmeFile, "## "+serverName)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // replace the contents of the dao/userExample.go file
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // replace the contents of the v1/userExample.proto file
			Old: protoFileMark,
			New: codes[parser.CodeTypeProto],
		},
		{ // replace the contents of the proto.sh file
			Old: protoShellFileGRPCMark,
			New: protoShellGRPCMark,
		},
		{ // replace the contents of the scripts/proto.sh file
			Old: protoShellFileMark,
			New: protoShellServiceTmplCode,
		},
		{ // replace the contents of the service/userExample_client_test.go file
			Old: serviceFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeService]),
		},
		{ // replace the contents of the Dockerfile file
			Old: dockerFileMark,
			New: dockerFileGrpcCode,
		},
		{ // replace the contents of the Dockerfile_build file
			Old: dockerFileBuildMark,
			New: dockerFileBuildGrpcCode,
		},
		{ // replace the contents of the image-build.sh file
			Old: imageBuildFileMark,
			New: imageBuildFileGrpcCode,
		},
		{ // replace the contents of the image-build-local.sh file
			Old: imageBuildLocalFileMark,
			New: imageBuildLocalFileGrpcCode,
		},
		{ // replace the contents of the docker-compose.yml file
			Old: dockerComposeFileMark,
			New: dockerComposeFileGrpcCode,
		},
		{ // replace the contents of the *-deployment.yml file
			Old: k8sDeploymentFileMark,
			New: k8sDeploymentFileGrpcCode,
		},
		{ // replace the contents of the *-svc.yml file
			Old: k8sServiceFileMark,
			New: k8sServiceFileGrpcCode,
		},
		{ // replace the configuration of the *.yml file
			Old: appConfigFileMark,
			New: rpcServerConfigCode,
		},
		// replace github.com/i2dou/sponge/templates/sponge
		{
			Old: selfPackageName + "/" + r.GetSourcePath(),
			New: moduleName,
		},
		// replace directory name
		{
			Old: strings.Join([]string{"api", "userExample", "v1"}, gofile.GetPathDelimiter()),
			New: strings.Join([]string{"api", serverName, "v1"}, gofile.GetPathDelimiter()),
		},
		{
			Old: "github.com/i2dou/sponge",
			New: moduleName,
		},
		{
			Old: moduleName + "/pkg",
			New: "github.com/i2dou/sponge/pkg",
		},
		{
			Old: "api/userExample/v1",
			New: fmt.Sprintf("api/%s/v1", serverName),
		},
		{
			Old: "api.userExample.v1",
			New: fmt.Sprintf("api.%s.v1", serverName), // protobuf package no "-" signs allowed
		},
		{
			Old: "sponge api docs",
			New: serverName + " api docs",
		},
		{
			Old: "_userExampleNO       = 2",
			New: fmt.Sprintf("_userExampleNO       = %d", rand.Intn(100)),
		},
		{
			Old: "serverNameExample",
			New: serverName,
		},
		// docker image and k8s deployment script replacement
		{
			Old: "server-name-example",
			New: xstrings.ToKebabCase(serverName), // snake_case to kebab_case
		},
		// docker image and k8s deployment script replacement
		{
			Old: "project-name-example",
			New: projectName,
		},
		{
			Old: "repo-addr-example",
			New: repoAddr,
		},
		{
			Old: "image-repo-host",
			New: repoHost,
		},
		{
			Old: "_grpcExample",
			New: "",
		},
		{
			Old: "_mixExample",
			New: "",
		},
		{
			Old: "root:123456@(192.168.3.37:3306)/account",
			New: dbDSN,
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
