package handler

import (
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/annotation"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

var annotationParser gencore.AnnotationParser = gencore.NewTagParser()

type QueryAnnotationMethodHandler struct {
}

func (q *QueryAnnotationMethodHandler) HandleMethod(context *generator.Context, methodArgs *gencore.MethodArgs, writer io.Writer, tracker namer.ImportTracker) bool {
	queryAnnotation := &annotation.Query{}
	ok, err := annotationParser.Parse(queryAnnotation, methodArgs.Method)
	gencore.FatalfOnErr(err)
	if !ok {
		return false
	}
	//TODO
	return true
}
