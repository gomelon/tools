package handler

import (
	"fmt"
	"github.com/huandu/xstrings"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
	"path"
	"path/filepath"
	"testing"
)

func TestName1(t *testing.T) {
	_, universe, _, err := construct(testFile)
	if err != nil {
		t.Fatal(err)
	}
	personDaoType := universe.Type(types.Name{Package: "foo/bar", Name: "PersonDao"})
	for methodName, method := range personDaoType.Methods {
		rowType := getRowType(method.Signature.Results[0])
		builder := NewSqlBuilder("person", "`", personDaoType, methodName, method, rowType, xstrings.ToSnakeCase)
		builder.Build()
		fmt.Println(builder.projection)
	}
}

func getRowType(typ *types.Type) *types.Type {
	for {
		if typ.Kind == types.Pointer || typ.Kind == types.Slice || typ.Kind == types.Map {
			typ = typ.Elem
		} else {
			return typ
		}
	}
}

func construct(files ...file) (*parser.Builder, types.Universe, []*types.Type, error) {
	testNamer := namer.NewPublicNamer(0)
	builder := parser.New()
	for _, f := range files {
		err := builder.AddFileForTest(path.Dir(f.path), filepath.FromSlash(f.path), []byte(f.contents))
		if err != nil {
			return nil, nil, nil, err
		}
	}
	universe, err := builder.FindTypes()
	if err != nil {
		return nil, nil, nil, err
	}
	orderer := namer.Orderer{Namer: testNamer}
	orderUniverse := orderer.OrderUniverse(universe)
	return builder, universe, orderUniverse, nil
}

var testFile = file{
	path: "foo/bar/common.go",
	contents: `
package bar

type Person struct {
	Id        int64
	Firstname string
	Lastname  string
	Age       int32
	Address   Address
}

type Address struct {
	Email    string
	location string
}

type NamesOnly struct {
	Firstname string
	Lastname  string
}

type PersonDao interface {
	//CountByLastname(lastname string) int32
	FindByLastname(lastname string) []*Person
	FindByFirstName(firstName string) []*NamesOnly
	FindFirstByOrderByLastnameAsc() *Person
	FindTopByOrderByAgeDesc() *Person
	QueryFirst10ByLastname(lastname string) []*Person
	FindTop3ByLastname(lastname string) []*Person
	//DeleteByLastname(lastname string) int32
	//RemoveByLastname(lastname string) int32
}
`,
}

type file struct {
	path     string
	contents string
}

type Person struct {
	Id        int64
	Firstname string
	Lastname  string
	Age       int32
	Address   Address
}

type Address struct {
	Email    string
	location string
}

type NamesOnly struct {
	Firstname string
	Lastname  string
}

type PersonDao interface {
	//CountByLastname(lastname string) int32
	FindByLastname(lastname string) []*Person
	FindByFirstName(firstName string) []*NamesOnly
	FindFirstByOrderByLastnameAsc() *Person
	FindTopByOrderByAgeDesc() *Person
	QueryFirst10ByLastname(lastname string) []*Person
	FindTop3ByLastname(lastname string) []*Person
	//DeleteByLastname(lastname string) int32
	//RemoveByLastname(lastname string) int32
}
