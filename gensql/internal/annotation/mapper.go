package annotation

import "k8s.io/gengo/types"

// NameMapper
// +melon.sql.Mapper
const NameMapper = "Mapper"

type Mapper struct {
	TableName string
}

func (g *Mapper) Kinds() []types.Kind {
	return []types.Kind{types.Interface}
}

func (g *Mapper) Namespace() string {
	return Namespace
}

func (g *Mapper) Name() string {
	return NameMapper
}
