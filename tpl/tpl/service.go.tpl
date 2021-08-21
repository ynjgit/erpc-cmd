package main 

import (
    "context"
    
    rpc "{{.GoPkg}}"
)

type {{.SvcNameLower}}Impl struct{}

{{range .Rpcs}}
// {{.Name}} {{.Comment}}
func (s *{{$.SvcNameLower}}Impl) {{.Name}}(ctx context.Context, req *rpc.{{.Req}}, rsp *rpc.{{.Rsp}}) error {
    // your business code ...

    return nil
}
{{end}}