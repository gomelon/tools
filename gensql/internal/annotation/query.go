package annotation

import "k8s.io/gengo/types"

// NameQuery
// +melon.sql.Query `SQL:"" `
const NameQuery = "Query"

type Query struct {
	SQL string
}

func (q Query) Kinds() []types.Kind {
	return []types.Kind{types.Func}
}

func (q Query) Namespace() string {
	return Namespace
}

func (q Query) Name() string {
	return NameQuery
}
