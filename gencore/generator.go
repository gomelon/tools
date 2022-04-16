package gencore

import (
	"flag"
	"github.com/Masterminds/sprig"
	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"text/template"
)

var NameSystemPublic = "public"
var NameSystemPrivate = "private"
var NameSystemRaw = "raw"

var DefaultNameSystems = namer.NameSystems{
	NameSystemPublic:  namer.NewPublicNamer(0),
	NameSystemPrivate: namer.NewPrivateNamer(0),
	NameSystemRaw:     namer.NewRawNamer("", nil),
}

func Args() *args.GeneratorArgs {
	arguments := args.Default()
	arguments.AddFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	curDir, err := os.Getwd()
	klog.Infoln("Current Dir:", curDir)
	if err != nil {
		klog.Fatalf("Error: %v", err)
	}
	//TODO 这里要增加自动设置InputDirs的逻辑
	klog.Infoln("Input Dir:", arguments.InputDirs)
	return arguments
}

func Packages() (types.Universe, error) {
	klog.Infoln("Start ...")
	arguments := Args()

	parserBuilder, err := arguments.NewBuilder()
	if err != nil {
		klog.Fatalf("Error: %v", err)
	}
	return parserBuilder.FindTypes()
}

func NewTemplate(context *generator.Context, name string) *template.Template {
	tmpl := template.New(name).
		Funcs(sprig.GenericFuncMap()).
		Funcs(template.FuncMap{
			"short": func(name string) string {
				return strings.ToLower(string(name[0]))
			},
		})

	for namerName, namer := range context.Namers {
		tmpl.Funcs(template.FuncMap{
			namerName: namer.Name,
		})
	}
	return tmpl
}

type TypeArgs struct {
	Type *types.Type
}

func NewTypeArgs(typ *types.Type) *TypeArgs {
	return &TypeArgs{
		Type: typ,
	}
}

type MethodArgs struct {
	TypeName      string
	TypeShortName string
	Type          *types.Type
	MethodName    string
	Method        *types.Type
	ParamNames    []string
	Params        []*types.Type
	ResultNames   []string
	Results       []*types.Type
}

func NewMethodArgs(typeName string, typ *types.Type, methodName string, methodType *types.Type) *MethodArgs {
	return &MethodArgs{
		TypeName:      typeName,
		TypeShortName: strings.ToLower(typeName[0:1]),
		Type:          typ,
		MethodName:    methodName,
		Method:        methodType,
		ParamNames:    methodType.Signature.ParameterNames,
		Params:        methodType.Signature.Parameters,
		ResultNames:   methodType.Signature.ResultNames,
		Results:       methodType.Signature.Results,
	}
}

func KindSliceToSet(annotation Annotation) map[types.Kind]bool {
	kindSet := make(map[types.Kind]bool)
	for _, kind := range annotation.Kinds() {
		kindSet[kind] = true
	}
	return kindSet
}
