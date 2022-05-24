func ({{.TypeShortName}} *{{.TypeName}}) {{.MethodName}}(
{{range $i, $v := .Params}}{{index $.ParamNames $i}} {{$v|raw}},{{end}}
)({{range $i, $v := .Results}}{{index $.ResultNames $i}} {{$v|raw}},{{end}}){
query := "{{.Extra.SQL}}"
db := melon.GetSqlExecutor(ctx, melon.DBNameDefault)
rows := melon.Must2(
db.QueryContext({{index $.ParamNames 0}}, query,
{{range $i, $v := .Params}}{{if $i}} sql.Named("{{index $.ParamNames $i}}", {{index $.ParamNames $i}}), {{end}} {{end}}
),
)
defer melon.Must(rows.Close())
if rows.Next() {
var user entity.User
melon.Must(rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Birthday, &user.CreatedAt))
return &user
} else {
return nil
}
}
