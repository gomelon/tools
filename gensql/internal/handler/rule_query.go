package handler

import (
	"fmt"
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/annotation"
	"github.com/gomelon/tools/gensql/internal/templates"
	"github.com/huandu/xstrings"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	"regexp"
	"strings"
)

type QueryRuleHandlerFactory struct {
}

func (q QueryRuleHandlerFactory) Create(genCtx *generator.Context, methodArgs *gencore.MethodArgs,
	writer io.Writer, tracker namer.ImportTracker) Handler {
	return &queryRuleHandler{
		genCtx:        genCtx,
		methodArgs:    methodArgs,
		writer:        writer,
		importTracker: tracker,
	}
}

type queryRuleHandler struct {
	genCtx        *generator.Context
	methodArgs    *gencore.MethodArgs
	writer        io.Writer
	importTracker namer.ImportTracker
}

var (
	queryPattern    = "Find|Read|Get|Query|Search|Stream"
	prefixTemplate  = regexp.MustCompile("^(" + queryPattern + ")((\\p{Lu}.*?))??By")
	keywordTemplate = "(%s)(?=(\\p{Lu}|\\P{InBASIC_LATIN}))"
)

func (q *queryRuleHandler) HandleRule() bool {

	if !q.shouldMatchRule() {
		return false
	}

	q.checkMethodResult()
	q.checkMethodParam()

	tplFileName, rowType := TplFileNameAndRowType(q.methodArgs)
	sql := q.generateSql(rowType)
	q.methodArgs.Extra["SQL"] = sql
	//TODO 使用该模版则表示需要melon,而melon需要注册数据库,包已经有在mod了,go在格式化会在该文件自动导入该包
	// 后续需要修改成手动去导入,否则很容易有问题
	tplBytes := gencore.Must(templates.FS.ReadFile(tplFileName + ".tpl"))
	tplText := string(tplBytes)
	tmpl := gencore.Must(gencore.NewTemplate(q.genCtx, "gen").Parse(tplText))
	gencore.Must1(tmpl.Execute(q.writer, q.methodArgs))
	return true
}

func (q *queryRuleHandler) shouldMatchRule() bool {
	methodName := q.methodArgs.MethodName
	return prefixTemplate.MatchString(methodName)
}

func (q *queryRuleHandler) checkMethodParam() {
	params := q.methodArgs.Params
	if len(params) == 0 {
		klog.Fatalf("invalid method signature: query method [%s.%s] must at least one parameter",
			q.methodArgs.Type.Name.Name, q.methodArgs.MethodName)
	}
}

func (q *queryRuleHandler) checkMethodResult() {
	results := q.methodArgs.Results
	resultLen := len(results)
	switch resultLen {
	case 0:
		gencore.FatalfOnErr(
			fmt.Errorf("invalid method signature: query method [%s.%s] must at least one result",
				q.methodArgs.Type.Name.Name, q.methodArgs.MethodName),
		)
	case 1:
		break
	case 2:
		if !gencore.IsError(results[resultLen-1]) {
			gencore.FatalfOnErr(
				fmt.Errorf("invalid method signature: query method [%s.%s]'s second result must is [error]",
					q.methodArgs.Type.Name.Name, q.methodArgs.MethodName),
			)
		}
	default:
		gencore.FatalfOnErr(
			fmt.Errorf("invalid method signature: query method [%s.%s] must not more than 2 result",
				q.methodArgs.Type.Name.Name, q.methodArgs.MethodName),
		)
	}
}

func (q *queryRuleHandler) generateSql(rowType *types.Type) string {
	var sql strings.Builder
	sql.WriteString("SELECT ")
	sql.WriteString(q.rowColumnNames(rowType))
	sql.WriteString(" ")
	sql.WriteString("FROM ")
	sql.WriteString(q.tableName(rowType))
	return sql.String()
}

func (q *queryRuleHandler) rowColumnNames(rowType *types.Type) string {
	columnNames := make([]string, 0, len(rowType.Members))
	for _, member := range rowType.Members {
		//TODO 这里要根据各方言转义字段
		//TODO 这里转成字段名要可以配置
		columnNames = append(columnNames, xstrings.ToSnakeCase(member.Name))
	}
	return strings.Join(columnNames, ", ")
}

func (q *queryRuleHandler) conditions(rowType *types.Type) string {
	methodName := q.methodArgs.MethodName
	return methodName
}

func (q *queryRuleHandler) tableName(rowType *types.Type) string {
	return q.mapperAnnotation().TableName
}

func (q *queryRuleHandler) rowFieldNames(rowType *types.Type) []string {
	fieldNames := make([]string, 0, len(rowType.Members))
	for _, member := range rowType.Members {
		if strings.ToTitle(member.Name) == member.Name {
			fieldNames = append(fieldNames, member.Name)
		}
	}
	return fieldNames
}

func (q *queryRuleHandler) mapperAnnotation() *annotation.Mapper {
	mapperAnnotation := q.methodArgs.TypeAnnotations[annotation.FullNameMapper]
	mapper := mapperAnnotation.(*annotation.Mapper)
	return mapper
}

//func parseCondition(methodName string) []any {
//	conditionsPart := prefixTemplate.ReplaceAllString(methodName, "")
//	strings.Index
//	return nil
//}

type condition struct {
	Operator   string
	ColumnName string
	Values     []any
}
