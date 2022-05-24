package handler

import (
	"github.com/gomelon/tools/gencore"
	"github.com/gomelon/tools/gensql/internal/annotation"
	"github.com/gomelon/tools/gensql/internal/sql"
	"github.com/gomelon/tools/gensql/internal/templates"
	"github.com/huandu/xstrings"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
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
	tplFileName, rowType := TplFileNameAndRowType(methodArgs)
	methodArgs.Extra["RowType"] = rowType

	querySQL := queryAnnotation.SQL
	methodArgs.Extra["SQL"] = querySQL

	parser, err := sql.CreateParser("", querySQL)
	gencore.FatalfOnErr(err)
	columns, err := parser.SelectColumns()
	gencore.FatalfOnErr(err)

	columnToMember := q.columnToMember(rowType)

	selectMembers := make([]types.Member, 0, len(columns))
	for _, column := range columns {
		member, ok := columnToMember[column.Alias]
		if !ok {
			klog.Fatalf("error: unknown column: method=[%s.%s],column=[%s]",
				methodArgs.Type.Name, methodArgs.MethodName, column.Alias)
		}
		selectMembers = append(selectMembers, member)
	}
	methodArgs.Extra["SelectMembers"] = selectMembers

	tplBytes := gencore.Must(templates.FS.ReadFile(tplFileName + ".tpl"))
	tplText := string(tplBytes)
	tmpl := gencore.Must(gencore.NewTemplate(context, "gen").Parse(tplText))
	gencore.Must1(tmpl.Execute(writer, methodArgs))
	return true
}

func (q *QueryAnnotationMethodHandler) columnToMember(rowType *types.Type) map[string]types.Member {
	columnToMember := make(map[string]types.Member, len(rowType.Members))
	for _, member := range rowType.Members {
		columnToMember[xstrings.ToSnakeCase(member.Name)] = member
	}
	return columnToMember
}
