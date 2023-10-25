package router

import (
	"text/template"
)

func init() {
	var err error
	handlerTmpl, err = template.New("iRouter").Parse(handlerTmplRaw)
	if err != nil {
		panic(err)
	}
}

var (
	handlerTmpl    *template.Template
	handlerTmplRaw = `
type {{$.Name}}Logicer interface {
{{range .MethodSet}}{{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

type {{$.Name}}Option func(*{{$.LowerName}}Options)

type {{$.LowerName}}Options struct {
	isFromRPC bool
	responser errcode.Responser
	zapLog    *zap.Logger
	httpErrors []*errcode.Error
	rpcStatus  []*errcode.RPCStatus
	wrapCtxFn  func(c *gin.Context) context.Context
}

func (o *{{$.LowerName}}Options) apply(opts ...{{$.Name}}Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func With{{$.Name}}HTTPResponse() {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.isFromRPC = false
	}
}

func With{{$.Name}}RPCResponse() {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.isFromRPC = true
	}
}

func With{{$.Name}}Responser(responser errcode.Responser) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.responser = responser
	}
}

func With{{$.Name}}Logger(zapLog *zap.Logger) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.zapLog = zapLog
	}
}

func With{{$.Name}}ErrorToHTTPCode(e ...*errcode.Error) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.httpErrors = e
	}
}

func With{{$.Name}}RPCStatusToHTTPCode(s ...*errcode.RPCStatus) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.rpcStatus = s
	}
}

func With{{$.Name}}WrapCtx(wrapCtxFn func(c *gin.Context) context.Context) {{$.Name}}Option {
	return func(o *{{$.LowerName}}Options) {
		o.wrapCtxFn = wrapCtxFn
	}
}

func Register{{$.Name}}Router(
	iRouter gin.IRouter,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iLogic {{$.Name}}Logicer,
	opts ...{{$.Name}}Option) {

	o := &{{$.LowerName}}Options{}
	o.apply(opts...)

	if o.responser == nil {
		o.responser = errcode.NewResponser(o.isFromRPC, o.httpErrors, o.rpcStatus)
	}
	if o.zapLog == nil {
		o.zapLog,_ = zap.NewProduction()
	}

	r := &{{$.LowerName}}Router {
		iRouter:               iRouter,
		groupPathMiddlewares:  groupPathMiddlewares,
		singlePathMiddlewares: singlePathMiddlewares,
		iLogic:                iLogic,
		iResponse:             o.responser,
		zapLog:                o.zapLog,
		wrapCtxFn:             o.wrapCtxFn,
	}
	r.register()
}

type {{$.LowerName}}Router struct {
	iRouter               gin.IRouter
	groupPathMiddlewares  map[string][]gin.HandlerFunc
	singlePathMiddlewares map[string][]gin.HandlerFunc
	iLogic                {{$.Name}}Logicer
	iResponse             errcode.Responser
	zapLog                *zap.Logger
	wrapCtxFn             func(c *gin.Context) context.Context
}

func (r *{{$.LowerName}}Router) register() {
{{range .Methods}}r.iRouter.Handle("{{.Method}}", "{{.Path}}", r.withMiddleware("{{.Method}}", "{{.Path}}", r.{{ .HandlerName }})...)
{{end}}
}

func (r *{{$.LowerName}}Router) withMiddleware(method string, path string, fn gin.HandlerFunc) []gin.HandlerFunc {
	handlerFns := []gin.HandlerFunc{}

	// determine if a route group is hit or miss, left prefix rule
	for groupPath, fns := range r.groupPathMiddlewares {
		if groupPath == "" || groupPath == "/" {
			handlerFns = append(handlerFns, fns...)
			continue
		}
		size := len(groupPath)
		if len(path) < size {
			continue
		}
		if groupPath == path[:size] {
			handlerFns = append(handlerFns, fns...)
		}
	}

	// determine if a single route has been hit
	key := strings.ToUpper(method) + "->" + path
	if fns, ok := r.singlePathMiddlewares[key]; ok {
		handlerFns = append(handlerFns, fns...)
	}

	return append(handlerFns, fn)
}

{{range .Methods}}
func (r *{{$.LowerName}}Router) {{ .HandlerName }} (c *gin.Context) {
	req := &{{.Request}}{}
	var err error
{{if .HasPathParams }}
	if err = c.ShouldBindUri(req); err != nil {
		r.zapLog.Warn("ShouldBindUri error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err = c.ShouldBindQuery(req); err != nil {
		r.zapLog.Warn("ShouldBindQuery error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err = c.ShouldBindJSON(req); err != nil {
		r.zapLog.Warn("ShouldBindJSON error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{else}}
	if err = c.ShouldBind(req); err != nil {
		r.zapLog.Warn("ShouldBind error", zap.Error(err), middleware.GCtxRequestIDField(c))
		r.iResponse.ParamError(c, err)
		return
	}
{{end}}
	var ctx context.Context
	if r.wrapCtxFn != nil {
		ctx = r.wrapCtxFn(c)
	} else {
		ctx = c
	}

	out, err := r.iLogic.{{.Name}}(ctx, req)
	if err != nil {
		r.iResponse.Error(c, err)
		return
	}

	r.iResponse.Success(c, out)
}
{{end}}
`
)
