package rule

import (
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/templates"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"strings"
)

type QueryRuleHandler struct {
}

var queryMethodPrefixNames = []string{"FindBy", "GetBy", "QueryBy", "SearchBy"}

func (q *QueryRuleHandler) HandleRule(context *generator.Context, methodArgs *gencore.MethodArgs,
	writer io.Writer, importTracker namer.ImportTracker) bool {
	if !q.ShouldMatchRule(methodArgs) {
		return false
	}

	//TODO 使用该模版则表示需要melon,而melon需要注册数据库,也以包已经有在mod了,go在格式化会在该文件自动导入该包
	tplBytes, err := templates.FS.ReadFile("GetWithCtx.tpl")
	gencore.FatalfOnErr(err)
	tplText := string(tplBytes)
	tmpl, err := gencore.NewTemplate(context, "gen").
		Parse(tplText)
	gencore.FatalfOnErr(err)
	err = tmpl.Execute(writer, methodArgs)
	gencore.FatalfOnErr(err)
	return true
}

func (q *QueryRuleHandler) ShouldMatchRule(methodArgs *gencore.MethodArgs) bool {
	methodName := methodArgs.MethodName
	for _, queryMethodPrefixName := range queryMethodPrefixNames {
		if strings.HasPrefix(methodName, queryMethodPrefixName) {
			return true
		}
	}
	return false
}
