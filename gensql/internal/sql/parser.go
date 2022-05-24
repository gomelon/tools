package sql

type Type int

const (
	TypeSelect Type = iota
	TypeInsert
	TypeUpdate
	TypeDelete
)

type Parser interface {
	Type() (Type, error)
	SelectColumns() ([]*Column, error)
}

func CreateParser(engin string, sql string) (Parser, error) {
	return NewMySQLParser(sql)
}

type Column struct {
	Alias          string
	TableQualifier string
}
