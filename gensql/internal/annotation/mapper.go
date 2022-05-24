package annotation

import (
	"k8s.io/gengo/types"
)

// NameMapper
// +melon.sql.Mapper
const (
	NameMapper     = "Mapper"
	FullNameMapper = Namespace + "." + NameMapper
)

type Mapper struct {
	TableName    string
	ColumnNaming string `default:"SnakeCase"` //SnakeCase: first_name CamelCase: SomeWords KebabCase: first-name
	Err          string `default:"db"`
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

func (g *Mapper) FullName() string {
	return FullNameMapper
}
