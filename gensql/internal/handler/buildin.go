package handler

import (
	"github.com/gomelon/tools/gencore"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

type BuildInMethodHandler struct {
}

func (b *BuildInMethodHandler) HandleMethod(context *generator.Context, methodArgs *gencore.MethodArgs, writer io.Writer, tracker namer.ImportTracker) bool {
	//TODO implement me
	return false
}
