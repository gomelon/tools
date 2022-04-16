func ({{.TypeShortName}} *{{.TypeName}}) {{.MethodName}}(
{{range $i, $v := .Params}}{{index $.ParamNames $i}} {{$v|raw}},{{end}})({{range $i, $v := .Results}}{{index $.ResultNames $i}} {{$v|raw}},{{end}}){
query := "SELECT `id`, `name`, `gender`, `birthday`, `created_at` " +
"FROM `user` " +
"WHERE `id` = ?"
db := melon.GetSqlExecutor(ctx, melon.DBNameDefault)
rows, err := db.Query(query, id)
melon.PanicOnError(err)
defer rows.Close()
if rows.Next() {
var user entity.User
err := rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Birthday, &user.CreatedAt)
melon.PanicOnError(err)
return &user
} else {
return nil
}
}
