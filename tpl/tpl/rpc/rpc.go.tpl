package {{.Pkg}}

import (
	"context"

	"github.com/ynjgit/erpc-go/client"
	"github.com/ynjgit/erpc-go/server"
)

{{if .Comment -}}
// {{.SvcName}} {{.Comment}}
{{- end}}
type {{.SvcName}} interface {
    {{range .Rpcs}}
    // {{.Name}} {{.Comment}}
    {{.Name}}(ctx context.Context, req *{{.Req}}, rsp *{{.Rsp}}) error
    {{end}}
}

{{range .Rpcs}}
func wrapperHandle{{.Name}}(ctx context.Context, svcImpl interface{}, reqBody []byte, f server.InterceptorFunc) (interface{}, error) {
    req := &{{.Req}}{}
    rsp := &{{.Rsp}}{}

    chain, err := f(reqBody, req)
	if err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, req interface{}, rsp interface{}) error {
		return svcImpl.({{$.SvcName}}).{{.Name}}(ctx, req.(*{{.Req}}), rsp.(*{{.Rsp}}))
	}

	err = chain.Handle(ctx, req, rsp, handler)
	return rsp, err
}
{{end}}

var {{.SvcName}}Methods = []server.ServiceMethod{
    {{range .Rpcs}}
    server.ServiceMethod{
        RPCName: "{{.FullName}}",
        Handle: wrapperHandle{{.Name}},
    },
    {{end}}
}

func Register{{.SvcName}}RPC(s *server.Server, svcImpl interface{}) {
    for _, svcMethod := range {{.SvcName}}Methods {
        svcMethod.SvcImpl = svcImpl
        s.AddServiceMethod(svcMethod)
    }    
}

type {{.SvcName}}Client interface {
    {{range .Rpcs}}
    // {{.Name}} {{.Comment}}
    {{.Name}}(ctx context.Context, req *{{.Req}}, opts ...client.Option) (*{{.Rsp}}, error)
    {{end}}
}

func New{{.SvcName}}Client(opts ...client.Option) {{.SvcName}}Client {
    return &{{.SvcNameLower}}ClientImpl{
		opts: opts,
	}
}

type {{.SvcNameLower}}ClientImpl struct {
	opts []client.Option
}

{{range .Rpcs}}
func (c *{{$.SvcNameLower}}ClientImpl) {{.Name}}(ctx context.Context, req *{{.Req}}, opts ...client.Option) (*{{.Rsp}}, error) {
    rsp := &{{.Rsp}}{}
    callopts := make([]client.Option, 0, len(c.opts)+len(opts)+2)
    callopts = append(callopts, c.opts...)
    callopts = append(callopts, client.WithSVCName("{{$.FullSvcName}}"))
    callopts = append(callopts, client.WithRPCName("{{.FullName}}"))
    callopts = append(callopts, opts...)
    err := client.DefaultClient.Call(ctx, req, rsp, callopts...)
    return rsp, err
}
{{end}}