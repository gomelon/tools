func ({{.TypeShortName}} *{{.TypeName}}) {{.MethodName}}({{range $i, $v := .Params}}{{index $.ParamNames $i}} {{$v|raw}},{{end}}
)({{range $i, $v := .Results}}{{index $.ResultNames $i}} {{$v|raw}},{{end}}){

query := "{{.Extra.SQL}}"
db := melon.GetSqlExecutor(ctx, melon.DBNameDefault)
rows, err := db.QueryContext({{index $.ParamNames 0}}, query,
{{range $i, $v := .Params}}{{if $i}} sql.Named("{{index $.ParamNames $i}}", {{index $.ParamNames $i}}), {{end}} {{end}}
)
if err != nil {
return nil, err
}
defer rows.Close()
var e {{.Extra.RowType|raw}}
if rows.Next() {
e = {{.Extra.RowType|raw}}{}
err = rows.Scan({{range $i, $v := .Extra.SelectMembers}}&e.{{$v.Name}},{{end}})
if err != nil {
return &e, err
}
}
return &e, rows.Err()
}
