package env

import (
	"fmt"
	"strings"

	"github.com/ynjgit/erpc-cmd/util/execute"
)

// Env the os env
type Env struct {
	HasProtoc      bool
	GoMajorVersion string
}

// GetEnv get the os env
func GetEnv() *Env {
	return defaulEnv
}

var defaulEnv = &Env{}

func init() {
	err := execute.HasCmd("protoc")
	if err == nil {
		defaulEnv.HasProtoc = true
	}

	// get go major version
	shellCmd := "go version | { read _ _ v _; echo ${v#go}; }"
	out, err := execute.RunCmd("sh", "-c", shellCmd)
	if err != nil {
		panic(err)
	}
	versions := strings.Split(out, ".")
	defaulEnv.GoMajorVersion = strings.Join(versions[0:len(versions)-1], ".")

	fmt.Println("env:", defaulEnv)
}
