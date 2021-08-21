package tpl

import (
	"os"
	"path/filepath"
	"text/template"
)

// GenMainGO gen the server go file: main.go
func GenMainGO(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "main.go.tpl")
	return excuteTpl(tpl, outfile, data)
}

// GenGOMod gen the server file: go.mod
func GenGOMod(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "go.mod.tpl")
	return excuteTpl(tpl, outfile, data)
}

// GenYaml gen the server config yaml: erpc_go.yaml
func GenYaml(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "erpc_go.yaml.tpl")
	return excuteTpl(tpl, outfile, data)
}

// GenRPCGO ...
func GenRPCGO(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "rpc/rpc.go.tpl")
	return excuteTpl(tpl, outfile, data)
}

// GenRPCGOMod ...
func GenRPCGOMod(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "rpc/go.mod.tpl")
	return excuteTpl(tpl, outfile, data)

}

// GenServiceGO gen the server  service file
func GenServiceGO(tplDir string, outfile string, data interface{}) error {
	tpl := filepath.Join(tplDir, "service.go.tpl")
	return excuteTpl(tpl, outfile, data)
}

func excuteTpl(tplfile string, outfile string, data interface{}) error {
	mainTmpl, err := template.ParseFiles(tplfile)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = mainTmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
