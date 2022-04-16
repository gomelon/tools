package main

import (
	"fmt"
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/annotation"
	"github.com/gomelon/tools/gensql/internal/handler"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	"path/filepath"
	"strings"
)

var annotationParser gencore.AnnotationParser = gencore.NewTagParser()

func main() {
	klog.InitFlags(nil)
	arguments := args.Default()
	if err := arguments.Execute(
		gencore.DefaultNameSystems,
		gencore.NameSystemPublic,
		Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("//go:build !%s\n// +build !%s\n\n",
		arguments.GeneratedBuildTag, arguments.GeneratedBuildTag)))
	for _, input := range context.Inputs {
		pkg := context.Universe[input]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		// if is target project code
		if len(pkg.SourcePath) == 0 {
			continue
		}
		klog.Infoln("Scan Source Dir:", input)
		for _, typ := range pkg.Types {
			mapperAnnotation := &annotation.Mapper{}

			ok, err := annotationParser.Parse(mapperAnnotation, typ)
			gencore.FatalfOnErr(err)
			if !ok {
				continue
			}
			if !gencore.IsAnnotationSupport(mapperAnnotation, typ) {
				klog.Warningf("%s.%s don't support type[%s]",
					mapperAnnotation.Namespace(), mapperAnnotation.Name(), typ.Kind)
				continue
			}
			klog.Infoln("Find Type", typ.Name.Name)

			pkgPath := pkg.Path
			if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
				pkgPath = strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
			}

			packages = append(packages,
				&generator.DefaultPackage{
					PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
					PackagePath: pkgPath,
					HeaderText:  header,
					GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
						return []generator.Generator{handler.NewTypeGenerator(pkg, typ, mapperAnnotation)}
					},
					FilterFunc: func(c *generator.Context, t *types.Type) bool {
						return t.Name.Package == pkg.Path && t.Kind == types.Interface
					},
				})
		}
	}
	return packages
}
