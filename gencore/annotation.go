package gencore

import (
	"bytes"
	"fmt"
	"k8s.io/gengo/types"
	"reflect"
	"strings"
)

var _ AnnotationParser = &TagAnnotationParser{}

type Annotation interface {
	Kinds() []types.Kind
	Namespace() string
	Name() string
	FullName() string
}

type AnnotationParser interface {
	Parse(annotation Annotation, typ *types.Type) (bool, error)
}
type TagAnnotationParser struct {
}

func NewTagParser() *TagAnnotationParser {
	return &TagAnnotationParser{}
}

func (t *TagAnnotationParser) Parse(annotation Annotation, typ *types.Type) (ok bool, err error) {
	has, tagStr, err := t.parseTagAnnotation(annotation.Namespace(), annotation.Name(), typ.CommentLines)
	if err != nil {
		return false, err
	}

	if !has {
		return false, nil
	}

	tag := reflect.StructTag(tagStr)

	annotationType := reflect.TypeOf(annotation).Elem()
	annotationValue := reflect.ValueOf(annotation).Elem()
	for i := 0; i < annotationType.NumField(); i++ {
		field := annotationType.Field(i)
		fieldName := field.Name
		value, ok := tag.Lookup(fieldName)

		if !ok {
			value = field.Tag.Get("default")
		}

		if len(value) == 0 && field.Tag.Get("required") == "false" {
			return false, fmt.Errorf("parse annotation error: field [%s] of %s.%s lack of value",
				fieldName, annotation.Namespace(), annotation.Name())
		}

		fieldValue := annotationValue.FieldByName(fieldName)

		err := SetValueFromString(fieldValue, value)

		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// IsAnnotationSupport Whether the annotation supports the type
func IsAnnotationSupport(annotation Annotation, typ *types.Type) bool {
	for _, kind := range annotation.Kinds() {
		if typ.Kind == kind {
			return true
		}
	}
	return false
}

func (t *TagAnnotationParser) parseTagAnnotation(namespace, name string, comments []string) (has bool, tagStr string, err error) {
	var annotation = namespace + "." + name
	lineIndex := -1

	lineIndex = t.annotationLineIndex(annotation, comments)
	if lineIndex < 0 {
		return
	}

	has = true

	tagStartIndex := strings.Index(comments[lineIndex], "`")
	if tagStartIndex < 0 {
		return
	}

	tagBuf := bytes.NewBuffer(make([]byte, 0, 32))
	tagEndIndex := tagStartIndex + 1 + strings.Index(comments[lineIndex][tagStartIndex+1:], "`")
	if tagEndIndex > tagStartIndex {
		tagBuf.WriteString(comments[lineIndex][tagStartIndex+1 : tagEndIndex])
		tagStr = strings.TrimSpace(tagBuf.String())
		return
	}
	tagBuf.WriteString(comments[lineIndex][tagStartIndex+1:])

	for _, comment := range comments[lineIndex+1:] {
		tagEndIndex = strings.Index(comment, "`")
		tagBuf.WriteString("\n")
		if tagEndIndex < 0 {
			tagBuf.WriteString(comment)
		} else {
			tagBuf.WriteString(comment[:tagEndIndex])
			break
		}
	}

	if tagEndIndex < 0 {
		err = fmt.Errorf("parse annotation error: %s lack of tag end char", annotation)
	}
	tagStr = strings.TrimSpace(tagBuf.String())
	return
}

func (t *TagAnnotationParser) annotationLineIndex(annotation string, comments []string) int {
	for i, comment := range comments {
		if strings.HasPrefix(strings.TrimSpace(comment), annotation) {
			return i
		}
	}
	return -1
}
