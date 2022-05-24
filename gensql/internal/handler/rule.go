package handler

import (
	"github.com/gomelon/tools/gencore"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

type RuleMethodHandler struct {
}

func (r *RuleMethodHandler) HandleMethod(context *generator.Context, methodArgs *gencore.MethodArgs,
	writer io.Writer, importTracker namer.ImportTracker) bool {
	for _, handlerFactory := range HandlerFactories {
		handler := handlerFactory.Create(context, methodArgs, writer, importTracker)
		handled := handler.HandleRule()
		if handled {
			return true
		}
	}
	return false
}

type RuleMethodHandlerFactory interface {
	Create(genCtx *generator.Context, methodArgs *gencore.MethodArgs,
		writer io.Writer, tracker namer.ImportTracker) Handler
}

var HandlerFactories = []RuleMethodHandlerFactory{&QueryRuleHandlerFactory{}}

type Handler interface {
	HandleRule() bool
}
