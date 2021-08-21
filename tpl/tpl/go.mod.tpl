module {{.Module}}

go {{.GoMajorVersion}}

replace github.com/ynjgit/erpc-go => ../../erpc-go
replace {{.GoPkg}} => ./rpc/{{.GoPkg}}