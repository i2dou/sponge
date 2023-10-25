// Package router is to generate gin router code.
package router

import (
	"bytes"
	"strings"

	"github.com/i2dou/sponge/cmd/protoc-gen-go-gin/internal/parse"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	stringsPkg         = protogen.GoImportPath("strings")
	contextPkg         = protogen.GoImportPath("context")
	errcodePkg         = protogen.GoImportPath("github.com/i2dou/sponge/pkg/errcode")
	middlewarePkg      = protogen.GoImportPath("github.com/i2dou/sponge/pkg/gin/middleware")
	zapPkg             = protogen.GoImportPath("go.uber.org/zap")
	ginPkg             = protogen.GoImportPath("github.com/gin-gonic/gin")
	deprecationComment = "// Deprecated: Do not use."
)

// GenerateFile generates a *_router.pb.go file.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_router.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by https://github.com/i2dou/sponge, DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	g.P("// import packages: ", stringsPkg.Ident(" "), contextPkg.Ident(" "), errcodePkg.Ident(" "),
		middlewarePkg.Ident(" "), zapPkg.Ident(" "), ginPkg.Ident(" "))
	g.P()

	for _, s := range file.Services {
		genService(file, g, s)
	}
	return g
}

func genService(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	// HTTP Server.
	sd := &tmplField{
		Name:      s.GoName,
		LowerName: strings.ToLower(s.GoName[:1]) + s.GoName[1:],
		FullName:  string(s.Desc.FullName()),
		FilePath:  file.Desc.Path(),
	}

	for _, m := range s.Methods {
		sd.Methods = append(sd.Methods, parse.GetMethods(m)...)
	}

	g.P(sd.execute())
}

type tmplField struct {
	Name      string // Greeter
	LowerName string // greeter
	FullName  string // v1.Greeter
	FilePath  string // api/v1/demo.proto

	Methods   []*parse.RPCMethod
	MethodSet map[string]*parse.RPCMethod
}

func (s *tmplField) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = map[string]*parse.RPCMethod{}
		for _, m := range s.Methods {
			m := m
			s.MethodSet[m.Name] = m
		}
	}

	buf := new(bytes.Buffer)
	if err := handlerTmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()
}