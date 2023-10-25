// Package service is to generate template code, router code, and error code.
package service

import (
	"bytes"

	"github.com/i2dou/sponge/cmd/protoc-gen-go-gin/internal/parse"

	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateFiles generate service logic, router, error code files.
func GenerateFiles(file *protogen.File) ([]byte, []byte, []byte) {
	if len(file.Services) == 0 {
		return nil, nil, nil
	}

	pss := parse.GetServices(file)
	logicContent := genServiceLogicFile(pss)
	routerFileContent := genRouterFile(pss)
	errCodeFileContent := genErrCodeFile(pss)

	return logicContent, routerFileContent, errCodeFileContent
}

func genServiceLogicFile(fields []*parse.PbService) []byte {
	lf := &serviceLogicFields{PbServices: fields}
	return lf.execute()
}

func genRouterFile(fields []*parse.PbService) []byte {
	rf := &routerFields{PbServices: fields}
	return rf.execute()
}

func genErrCodeFile(fields []*parse.PbService) []byte {
	cf := &errCodeFields{PbServices: fields}
	return cf.execute()
}

type serviceLogicFields struct {
	PbServices []*parse.PbService
}

func (f *serviceLogicFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := serviceLogicTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return handleSplitLineMark(buf.Bytes())
}

type routerFields struct {
	PbServices []*parse.PbService
}

func (f *routerFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := routerTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	return handleSplitLineMark(buf.Bytes())
}

type errCodeFields struct {
	PbServices []*parse.PbService
}

func (f *errCodeFields) execute() []byte {
	buf := new(bytes.Buffer)
	if err := rpcErrCodeTmpl.Execute(buf, f); err != nil {
		panic(err)
	}
	data := bytes.ReplaceAll(buf.Bytes(), []byte("// --blank line--"), []byte{})
	return handleSplitLineMark(data)
}

var splitLineMark = []byte(`// ---------- Do not delete or move this split line, this is the merge code marker ----------`)

func handleSplitLineMark(data []byte) []byte {
	ss := bytes.Split(data, splitLineMark)
	if len(ss) <= 2 {
		return ss[0]
	}

	var out []byte
	for i, s := range ss {
		out = append(out, s...)
		if i < len(ss)-2 {
			out = append(out, splitLineMark...)
		}
	}
	return out
}
