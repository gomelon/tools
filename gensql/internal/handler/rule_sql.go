package handler

import (
	"fmt"
	"github.com/samber/lo"
	"k8s.io/gengo/types"
	"regexp"
	"strings"
)

type SqlBuilder struct {
	Table            string
	QuotedIdentifier string
	Type             *types.Type
	MethodName       string
	Method           *types.Type
	RowType          *types.Type
	FieldConverter   func(string) string
	methodNameParts  *methodNameParts
	fieldPaths       [][]string
	operatorType     OperatorType
	projection       string
	condition        string
	sort             string
}

type OperatorType int

const (
	OperatorTypeInvalid OperatorType = iota
	OperatorTypeQuery
	OperatorTypeCount
	OperatorTypeExists
	OperatorTypeUpdate
	OperatorTypeDelete
)

type SubjectType struct {
	ParamCount int
	Name       string
	Keywords   []string
}

var (
	SubjectTypeDistinct = &SubjectType{ParamCount: 0, Name: "Distinct", Keywords: []string{"Distinct"}}
)

type PredicationType struct {
	ParamCount int
	Name       string
	Keywords   []string
}

var (
	And        = &PredicationType{ParamCount: 0, Name: "And", Keywords: []string{"And"}}
	Or         = &PredicationType{ParamCount: 0, Name: "Or", Keywords: []string{"Or"}}
	After      = &PredicationType{ParamCount: 1, Name: "After", Keywords: []string{"After", "IsAfter"}}
	Before     = &PredicationType{ParamCount: 1, Name: "Before", Keywords: []string{"Before", "IsBefore"}}
	Between    = &PredicationType{ParamCount: 2, Name: "Between", Keywords: []string{"Between"}}
	Contains   = &PredicationType{ParamCount: 1, Name: "Containing", Keywords: []string{"Contains"}}
	Exists     = &PredicationType{ParamCount: 1, Name: "Exists", Keywords: []string{"Exists"}}
	True       = &PredicationType{ParamCount: 0, Name: "True", Keywords: []string{"True", "IsTrue"}}
	False      = &PredicationType{ParamCount: 0, Name: "False", Keywords: []string{"False", "IsFalse"}}
	Is         = &PredicationType{ParamCount: 1, Name: "Is", Keywords: []string{"Is", "Equals", ""}}
	Not        = &PredicationType{ParamCount: 1, Name: "Not", Keywords: []string{"Not", "IsNot"}}
	Gt         = &PredicationType{ParamCount: 1, Name: "Gt", Keywords: []string{"Gt"}}
	Gte        = &PredicationType{ParamCount: 1, Name: "Gte", Keywords: []string{"Gte"}}
	Lt         = &PredicationType{ParamCount: 1, Name: "Lt", Keywords: []string{"Lt"}}
	Lte        = &PredicationType{ParamCount: 1, Name: "Lte", Keywords: []string{"Lte"}}
	In         = &PredicationType{ParamCount: 1, Name: "In", Keywords: []string{"In", "IsIn"}}
	NotIn      = &PredicationType{ParamCount: 1, Name: "NotIn", Keywords: []string{"NotIn", "IsNotIn"}}
	IsEmpty    = &PredicationType{ParamCount: 0, Name: "IsEmpty", Keywords: []string{"IsEmpty", "Empty"}}
	IsNotEmpty = &PredicationType{ParamCount: 0, Name: "IsNotEmpty", Keywords: []string{"IsNotEmpty", "NotEmpty"}}
	IsNotNull  = &PredicationType{ParamCount: 0, Name: "IsNotNull", Keywords: []string{"NotNull", "IsNotNull"}}
	IsNull     = &PredicationType{ParamCount: 0, Name: "IsNull", Keywords: []string{"Null", "IsNull"}}
	Like       = &PredicationType{ParamCount: 1, Name: "Like", Keywords: []string{"Like", "IsLike"}}
	HasPrefix  = &PredicationType{ParamCount: 1, Name: "HasPrefix", Keywords: []string{"HasPrefix", "StartsWith"}}
	HasSuffix  = &PredicationType{ParamCount: 1, Name: "HasSuffix", Keywords: []string{"HasSuffix", "EndsWith"}}
	NotLike    = &PredicationType{ParamCount: 1, Name: "NotLike", Keywords: []string{"NotLike", "IsNotLike"}}
	Regex      = &PredicationType{ParamCount: 1, Name: "Regex", Keywords: []string{"Regex", "MatchesRegex", "Matches"}}
	// TODO 待实现
	//Near   = &PredicationType{ParamCount: 0, Name: "Near", Keywords: []string{"Near", "IsNear"}}
	//Within = &PredicationType{ParamCount: 0, Name: "Within", Keywords: []string{"Within", "IsWithin"}}
)

var PredicationTypes = []*PredicationType{
	And, Or, After, Before, Between, Contains, Exists, True,
	False, Is, Not, Gt, Gte, Lt, Lte, In,
	NotIn, IsEmpty, IsNotEmpty, IsNotNull, IsNull, Like, HasPrefix, HasSuffix,
	NotLike, Regex, //Near, Within,
}

func NewSqlBuilder(table string, quotedIdentifier string, typ *types.Type, methodName string, method *types.Type,
	rowType *types.Type, fieldConverter func(string) string) *SqlBuilder {
	return &SqlBuilder{
		Table:            table,
		QuotedIdentifier: quotedIdentifier,
		Type:             typ,
		MethodName:       methodName,
		Method:           method,
		RowType:          rowType,
		FieldConverter:   fieldConverter,
	}
}

func (b *SqlBuilder) Build() {
	b.SplitMethodName()
	b.buildFieldNames([]string{}, b.RowType.Members)
	b.buildOperatorType()
	b.buildProjection()
	b.buildCondition()
}

func (b *SqlBuilder) SplitMethodName() {
	seqPattern := regexp.MustCompile("(\\d+|_|([A-Z]+[a-z]*))")
	parts := seqPattern.FindAllStringIndex(b.MethodName, len(b.MethodName))
	methodNameParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if "_" == b.MethodName[part[0]:part[1]] {
			methodNameParts = append(methodNameParts, ".")
		} else {
			methodNameParts = append(methodNameParts, b.MethodName[part[0]:part[1]])
		}
	}
	b.methodNameParts = newMethodNameParts(methodNameParts)
}

func (b *SqlBuilder) buildOperatorType() {
	methodNamePart := b.methodNameParts.getAndNext()
	if lo.Contains([]string{"Get", "Find", "Query", "Search", "Read"}, methodNamePart) {
		b.operatorType = OperatorTypeQuery
	} else if methodNamePart == "Count" {
		b.operatorType = OperatorTypeCount
	} else if methodNamePart == "Exists" {
		b.operatorType = OperatorTypeExists
	} else if methodNamePart == "Update" {
		b.operatorType = OperatorTypeUpdate
	} else if methodNamePart == "Delete" {
		b.operatorType = OperatorTypeDelete
	}
}

func (b *SqlBuilder) buildFieldNames(parentMemberNames []string, members []types.Member) {
	memberDepth := len(parentMemberNames)
	for _, member := range members {
		memberNames := make([]string, memberDepth+1)
		memberNames = append(memberNames, parentMemberNames...)
		memberNames = append(memberNames, member.Name)
		if member.Type.Kind == types.Struct {
			b.buildFieldNames(memberNames, member.Type.Members)
		} else {
			b.fieldPaths = append(b.fieldPaths, memberNames)
		}
	}
}

func (b *SqlBuilder) buildProjection() {
	var projection *strings.Builder
	fieldPath, matched := b.tryMatchField()
	if matched {
		b.quoteField(projection, fieldPath)
		b.projection = projection.String()
		return
	}

	if b.methodNameParts.get() == "Distinct" {
		projection.WriteString("DISTINCT")
	}

	for i, fieldPath := range b.fieldPaths {
		if i > 0 {
			projection.WriteString(",")
		}
		b.quoteField(projection, fieldPath)
	}

	b.projection = projection.String()
}

func (b *SqlBuilder) buildCondition() {
	methodNamePart := b.methodNameParts.getAndNext()
	if methodNamePart != "By" {
		panic(
			fmt.Errorf(
				"invalid method signature: method [%s.%s]'s expect [By] but is [%s]",
				b.Type.Name.Name, b.MethodName, methodNamePart),
		)
	}
	var condition *strings.Builder
	fieldPath, matched := b.tryMatchField()
	if !matched {
		panic(
			fmt.Errorf(
				"invalid method signature: query method [%s.%s]'s expect a field name but is [%s]",
				b.Type.Name.Name, b.MethodName, methodNamePart),
		)
	}
	b.quoteField(condition, fieldPath)
	//methodNamePart = b.methodNameParts.getAndNext()
	//for _, predicationType := range PredicationTypes {
	//	predicationType.Name
	//}

}

func (b *SqlBuilder) quoteField(strBuilder *strings.Builder, fieldPath []string) {
	strBuilder.WriteString(b.QuotedIdentifier)
	b.convertToFieldName(strBuilder, fieldPath)
	strBuilder.WriteString(b.QuotedIdentifier)
}

func (b *SqlBuilder) convertToFieldName(strBuilder *strings.Builder, fieldPath []string) {
	for j, fieldName := range fieldPath {
		if j > 0 {
			strBuilder.WriteString("_")
		}
		strBuilder.WriteString(b.FieldConverter(fieldName))
	}
}

func (b *SqlBuilder) tryMatchField() ([]string, bool) {
	// multiple part assemble fieldName
	for _, fieldPath := range b.fieldPaths {
		if b.methodNameParts.tryMatch(fieldPath) {
			return fieldPath, true
		}
	}
	return []string{}, false
}

type methodNameParts struct {
	parts    []string
	position int
}

func newMethodNameParts(parts []string) *methodNameParts {
	return &methodNameParts{parts: parts, position: 0}
}

func (m *methodNameParts) getAndNext() string {
	part := m.parts[m.position]
	m.position++
	return part
}

func (m *methodNameParts) get() string {
	return m.parts[m.position]
}

func (m *methodNameParts) setPosition(position int) {
	m.position = position
}

func (m *methodNameParts) tryMatch(path []string) bool {
	nextIndex := 0
	pathPartIndex := 0
	partLen := len(m.parts)
	pathLen := len(path)
	for i := m.position; m.position < partLen && pathPartIndex < pathLen; i++ {
		pathPart := path[pathPartIndex]
		part := m.parts[m.position]
		index := strings.Index(pathPart, part)
		if nextIndex != index {
			break
		}
		nextIndex = nextIndex + len(part)
		if nextIndex == len(pathPart) {
			pathPartIndex++
		}
	}
	return pathPartIndex == pathLen
}

type Criteria struct {
	Projections []Projection
	Condition   []Condition
	Pagination  Pagination
}

type Projection struct {
	Operator    string
	Field       string
	Projections []*Projection
}

type Condition struct {
	Operator   string
	Field      string
	Params     []interface{}
	Conditions []Condition
}

type Pagination struct {
	Offset      int
	Limit       int
	SearchCount bool
}

type Sort struct {
	Field string
}
