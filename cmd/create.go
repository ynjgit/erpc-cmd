package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/ynjgit/erpc-cmd/env"
	"github.com/ynjgit/erpc-cmd/tpl"
	"github.com/ynjgit/erpc-cmd/util/execute"
)

var (
	protofile string
	outputDir string
	tplDir    string
)

type protoRPC struct {
	FullName string
	Name     string
	Comment  string
	Req      string
	Rsp      string
}

type protoServie struct {
	FullPkg      string
	Pkg          string
	GoPkg        string
	FullSvcName  string
	SvcName      string
	SvcNameLower string
	Comment      string
	Rpcs         map[string]*protoRPC
}

type protoDefine struct {
	GoMajorVersion string
	Module         string
	GoPkg          string
	Pkg            string
	FullPkg        string
	App            string
	Server         string
	Svcs           map[string]*protoServie
}

var (
	pd = &protoDefine{Svcs: make(map[string]*protoServie)}
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create the project with the input protobuf: *.proto",
	PreRunE: createPreRun,
	RunE:    create,
}

func create(cmd *cobra.Command, args []string) error {
	// parse protofile
	reader, _ := os.Open(protofile)
	defer reader.Close()

	parser := proto.NewParser(reader)
	pb, _ := parser.Parse()

	proto.Walk(pb, proto.WithOption(walkOption), proto.WithPackage(walkPackage), proto.WithService(walkService))
	fmt.Println("fullpkg", pd.FullPkg)
	pkg := strings.Split(pd.FullPkg, ".")
	if len(pkg) != 3 || pkg[0] != "erpc" {
		return fmt.Errorf("proto package must like erpc.{app}.{server}")
	}
	pd.App = pkg[1]
	pd.Server = pkg[2]
	pd.Pkg = pkg[2]
	pd.Module = pkg[2]

	fmt.Println("packege:", pd)
	fmt.Println("go_packege:", pd.GoPkg)
	if len(pd.Svcs) == 0 {
		return fmt.Errorf("no service got")
	}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(outputDir, "rpc"), os.ModePerm)
	if err != nil {
		return err
	}

	// protoc
	eMsg, err := execute.RunCmd(
		"protoc",
		fmt.Sprintf("--go_out=%s", filepath.Join(outputDir, "rpc")),
		fmt.Sprintf("--proto_path=%s", filepath.Dir(protofile)),
		protofile)
	if err != nil {
		return err
	}
	fmt.Println(eMsg)

	fmt.Println("tplDir:", tplDir)
	genfile := filepath.Join(outputDir, "main.go")
	err = tpl.GenMainGO(tplDir, genfile, pd)
	if err != nil {
		return err
	}

	genfile = filepath.Join(outputDir, "go.mod")
	err = tpl.GenGOMod(tplDir, genfile, pd)
	if err != nil {
		return err
	}

	genfile = filepath.Join(outputDir, "erpc_go.yaml")
	err = tpl.GenYaml(tplDir, genfile, pd)
	if err != nil {
		return err
	}

	rpcDir := filepath.Join(outputDir, "rpc", pd.GoPkg)
	err = os.MkdirAll(rpcDir, os.ModePerm)
	if err != nil {
		return err
	}

	genfile = filepath.Join(rpcDir, "go.mod")
	err = tpl.GenRPCGOMod(tplDir, genfile, pd)
	if err != nil {
		return err
	}

	for _, svc := range pd.Svcs {
		svc.FullPkg = pd.FullPkg
		svc.Pkg = pd.Pkg
		svc.GoPkg = pd.GoPkg
		genfile = filepath.Join(rpcDir, fmt.Sprintf("%s_rpc.go", strcase.ToSnake(svc.SvcName)))
		err := tpl.GenRPCGO(tplDir, genfile, svc)
		if err != nil {
			return err
		}

		genfile = filepath.Join(outputDir, fmt.Sprintf("%s.go", strcase.ToSnake(svc.SvcName)))
		err = tpl.GenServiceGO(tplDir, genfile, svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func walkOption(p *proto.Option) {
	if p.Name == "go_package" {
		pd.GoPkg = p.Constant.Source
	}
}

func walkPackage(p *proto.Package) {
	pd.FullPkg = p.Name
}

func walkService(s *proto.Service) {
	svcName := s.Name
	if _, ok := pd.Svcs[svcName]; !ok {
		svc := &protoServie{
			FullPkg:      pd.FullPkg,
			FullSvcName:  fmt.Sprintf("%s.%s", pd.FullPkg, s.Name),
			SvcName:      s.Name,
			SvcNameLower: strings.ToLower(s.Name),
			Rpcs:         make(map[string]*protoRPC),
		}
		if s.Comment != nil {
			svc.Comment = s.Comment.Message()
		}
		pd.Svcs[svcName] = svc
	}
	svc := pd.Svcs[svcName]

	for _, elem := range s.Elements {
		if rpc, ok := elem.(*proto.RPC); ok {
			handleRPC(svc, rpc)
		}
	}
}

func handleRPC(svc *protoServie, rpc *proto.RPC) {
	pbRPC := &protoRPC{
		FullName: fmt.Sprintf("/%s.%s/%s", svc.FullPkg, svc.SvcName, rpc.Name),
		Name:     rpc.Name,
		Req:      rpc.RequestType,
		Rsp:      rpc.ReturnsType,
	}
	if rpc.Comment != nil {
		pbRPC.Comment = rpc.Comment.Message()
	}

	svc.Rpcs[rpc.Name] = pbRPC
}

func createPreRun(cmd *cobra.Command, args []string) error {
	// check environment
	e := env.GetEnv()
	if !e.HasProtoc {
		return fmt.Errorf("protoc not found, please install... ")
	}

	// check flags
	flags := cmd.Flags()
	protofile, _ = flags.GetString("protofile")
	_, err := os.Stat(protofile)
	if os.IsNotExist(err) {
		return err
	}

	outputDir, _ = flags.GetString("outputdir")
	tplDir, _ = flags.GetString("tpldir")

	pd.GoMajorVersion = env.GetEnv().GoMajorVersion

	return nil
}

func init() {
	createCmd.Flags().StringP("protofile", "p", "", "the input protofile: *.proto")
	createCmd.MarkFlagRequired("protofile")
	createCmd.Flags().StringP("outputdir", "o", "./", "output dir to generate the server code")
	createCmd.Flags().StringP("tpldir", "t", "~/.erpc/tpl", "template dir")
}
