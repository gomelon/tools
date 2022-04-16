package handler

import (
	"github.com/gomelon/tools/gencore"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

type TypeGenerator struct {
	generator.DefaultGen
	Package       *types.Package
	Annotation    gencore.Annotation
	Typ           *types.Type
	ImportTracker namer.ImportTracker
}

func NewTypeGenerator(pkg *types.Package, typ *types.Type, annotation gencore.Annotation) TypeGenerator {
	return TypeGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: typ.Name.Name + "Impl",
		},
		Package:       pkg,
		Typ:           typ,
		Annotation:    annotation,
		ImportTracker: generator.NewImportTracker(),
	}
}

func (t TypeGenerator) Filter(context *generator.Context, typ *types.Type) bool {
	for _, r := range t.Annotation.Kinds() {
		if t.Typ.Kind == r {
			return true
		}
	}
	return false
}

func (t TypeGenerator) Namers(context *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		gencore.NameSystemRaw: namer.NewRawNamer(t.Package.Path, t.ImportTracker),
	}
}

func (t TypeGenerator) GenerateType(context *generator.Context, typ *types.Type, writer io.Writer) error {
	tmpl := gencore.NewTemplate(context, "gen")
	typeArgs := gencore.NewTypeArgs(typ)
	tmpl, err := tmpl.Parse(`
type {{.Type|raw}}Impl struct {
}
`)
	gencore.FatalfOnErr(err)
	err = tmpl.Execute(writer, typeArgs)
	gencore.FatalfOnErr(err)
	for methodName, methodType := range typ.Methods {
		methodArgs := gencore.NewMethodArgs(typ.Name.Name+"Impl", typ, methodName, methodType)
		GenerateMethod(context, methodArgs, writer, t.ImportTracker)
	}

	return nil
}

var methodHandlers = []MethodHandler{
	&QueryAnnotationMethodHandler{}, &BuildInMethodHandler{}, &RuleMethodHandler{},
}

func GenerateMethod(context *generator.Context, methodArgs *gencore.MethodArgs,
	writer io.Writer, importTracker namer.ImportTracker) {
	for _, methodHandler := range methodHandlers {
		generated := methodHandler.HandleMethod(context, methodArgs, writer, importTracker)
		if generated {
			break
		}
	}
}

type MethodHandler interface {
	HandleMethod(context *generator.Context, methodArgs *gencore.MethodArgs,
		writer io.Writer, importTracker namer.ImportTracker) bool
}
