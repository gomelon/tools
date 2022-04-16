package handler

import (
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/handler/rule"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

type RuleMethodHandler struct {
}

func (r *RuleMethodHandler) HandleMethod(context *generator.Context, methodArgs *gencore.MethodArgs,
	writer io.Writer, importTracker namer.ImportTracker) bool {
	for _, handler := range RuleHandlers {
		handled := handler.HandleRule(context, methodArgs, writer, importTracker)
		if handled {
			return true
		}
	}
	return false
}

var RuleHandlers = []RuleHandler{&rule.QueryRuleHandler{}}

type RuleHandler interface {
	HandleRule(context *generator.Context, methodArgs *gencore.MethodArgs, writer io.Writer, tracker namer.ImportTracker) bool
}
