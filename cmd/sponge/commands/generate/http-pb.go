package generate

import (
	"errors"
	"fmt"

	"github.com/i2dou/sponge/pkg/replacer"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// HTTPPbCommand generate web service code based on protobuf file
func HTTPPbCommand() *cobra.Command {
	var (
		moduleName   string // module name for go.mod
		serverName   string // server name
		projectName  string // project name for deployment name
		repoAddr     string // image repo address
		outPath      string // output directory
		protobufFile string // protobuf file, support * matching
	)

	cmd := &cobra.Command{
		Use:   "http-pb",
		Short: "Generate web service code based on protobuf file",
		Long: `generate web service code based on protobuf file.

Examples:
  # generate web service code.
  sponge web http-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./test.proto

  # generate web service code and specify the output directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge web http-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./test.proto --out=./yourServerDir

  # generate web service code and specify the docker image repository address.
  sponge web http-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --protobuf-file=./test.proto
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName, serverName = convertProjectAndServerName(projectName, serverName)
			return runGenHTTPPbCommand(moduleName, serverName, projectName, protobufFile, repoAddr, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&protobufFile, "protobuf-file", "f", "", "proto file")
	_ = cmd.MarkFlagRequired("protobuf-file")

	cmd.Flags().StringVarP(&repoAddr, "repo-addr", "r", "", "docker image repository address, excluding http and repository names")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./serverName_http-pb_<time>")

	return cmd
}

func runGenHTTPPbCommand(moduleName string, serverName string, projectName string, protobufFile string, repoAddr string, outPath string) error {
	protobufFiles, isImportTypes, err := parseProtobufFiles(protobufFile)
	if err != nil {
		return err
	}

	subTplName := "http-pb"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{ // processing-only subdirectories
		"api/types", "cmd/serverNameExample_httpPbExample",
		"sponge/configs", "sponge/deployments", "sponge/docs", "sponge/scripts", "sponge/third_party",
		"internal/config", "internal/ecode", "internal/routers", "internal/server",
	}
	subFiles := []string{ // processing of sub-documents only
		"sponge/.gitignore", "sponge/.golangci.yml", "sponge/go.mod", "sponge/go.sum",
		"sponge/Jenkinsfile", "sponge/Makefile", "sponge/README.md",
	}
	ignoreDirs := []string{} // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{ // specify the files in the subdirectory to be ignored for processing
		"types.pb.validate.go", "types.pb.go", // api/types
		"swagger.json", "swagger.yaml", "apis.swagger.json", "apis.html", "docs.go", // sponge/docs
		"userExample_rpc.go", "systemCode_rpc.go", "userExample_http.go", // internal/ecode
		"routers_pbExample_test.go", "routers.go", "routers_test.go", "userExample.go", "userExample_router.go", // internal/routers
		"grpc.go", "grpc_option.go", "grpc_test.go", // internal/server
	}

	if !isImportTypes {
		ignoreFiles = append(ignoreFiles, "types.proto")
	}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addHTTPPbFields(moduleName, serverName, projectName, repoAddr, r)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err = r.SaveFiles(); err != nil {
		return err
	}

	_ = saveProtobufFiles(moduleName, serverName, r.GetOutputDir(), protobufFiles)
	_ = saveGenInfo(moduleName, serverName, r.GetOutputDir())
	_ = saveEmptySwaggerJSON(r.GetOutputDir())

	fmt.Printf(`
using help:
  1. open a terminal and execute the command to generate code: make proto
  2. open file internal/handler/xxx.go, replace panic("implement me") according to template code example.
  3. compile and run service: make run
  4. visit http://localhost:8080/apis/swagger/index.html in your browser, and test api interface.

`)
	fmt.Printf("generate %s's web service code successfully, out = %s\n", serverName, r.GetOutputDir())
	return nil
}

func addHTTPPbFields(moduleName string, serverName string, projectName string, repoAddr string,
	r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	repoHost, _ := parseImageRepoAddr(repoAddr)

	fields = append(fields, deleteFieldsMark(r, httpFile, startMark, endMark)...)
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
		{ // replace the contents of the Dockerfile file
			Old: dockerFileMark,
			New: dockerFileHTTPCode,
		},
		{ // replace the contents of the Dockerfile_build file
			Old: dockerFileBuildMark,
			New: dockerFileBuildHTTPCode,
		},
		{ // replace the contents of the image-build.sh file
			Old: imageBuildFileMark,
			New: imageBuildFileHTTPCode,
		},
		{ // replace the contents of the image-build-local.sh file
			Old: imageBuildLocalFileMark,
			New: imageBuildLocalFileHTTPCode,
		},
		{ // replace the contents of the docker-compose.yml file
			Old: dockerComposeFileMark,
			New: dockerComposeFileHTTPCode,
		},
		{ // replace the contents of the *-deployment.yml file
			Old: k8sDeploymentFileMark,
			New: k8sDeploymentFileHTTPCode,
		},
		{ // replace the contents of the *-svc.yml file
			Old: k8sServiceFileMark,
			New: k8sServiceFileHTTPCode,
		},
		{ // replace the configuration of the *.yml file
			Old: appConfigFileMark,
			New: httpServerConfigCode,
		},
		{ // replace the contents of the proto.sh file
			Old: protoShellFileGRPCMark,
			New: "",
		},
		{ // replace the contents of the proto.sh file
			Old: protoShellFileMark,
			New: protoShellHandlerCode,
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
			Old: "sponge api docs",
			New: serverName + " api docs",
		},
		{
			Old: "serverNameExample",
			New: serverName,
		},
		// docker image and k8s deployment script replacement
		{
			Old: "server-name-example",
			New: xstrings.ToKebabCase(serverName), // convert to kebab-case format
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
			Old: "_httpPbExample",
			New: "",
		},
		{
			Old: "_pbExample",
			New: "",
		},
		{
			Old: "_mixExample",
			New: "",
		},
	}...)

	return fields
}
