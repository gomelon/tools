package gencore

import (
	"errors"
	"k8s.io/gengo/types"
	"reflect"
	"strconv"
)

func SetValueFromString(v reflect.Value, strVal string) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(strVal, 0, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(val) {
			return errors.New("Int value too big: " + strVal)
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(strVal, 0, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(val) {
			return errors.New("UInt value too big: " + strVal)
		}
		v.SetUint(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.Float64:
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.String:
		v.SetString(strVal)
	case reflect.Bool:
		val, err := strconv.ParseBool(strVal)
		if err != nil {
			return err
		}
		v.SetBool(val)
	default:
		return errors.New("Unsupported kind: " + v.Kind().String())
	}
	return nil
}

func Implements(subType *types.Type, ifc any) bool {
	ifcType := reflect.TypeOf(ifc).Elem()
	if ifcType.Kind() != reflect.Interface {
		return false
	}
	subTypeMethods := subType.Methods
	for i := 0; i < ifcType.NumMethod(); i++ {
		ifcMethod := ifcType.Method(i)
		methodName := ifcMethod.Name
		subTypeMethod, ok := subTypeMethods[methodName]
		if !ok {
			return false
		}
		subTypeParams := subTypeMethod.Signature.Parameters
		subTypeResults := subTypeMethod.Signature.Results

		ifcMethodType := ifcMethod.Type
		if len(subTypeParams) != ifcMethodType.NumIn() ||
			len(subTypeResults) != ifcMethodType.NumOut() {
			return false
		}
		for j := 0; j < ifcMethodType.NumIn(); j++ {
			ifcParam := ifcMethodType.In(j)
			subTypeParam := subTypeParams[j]
			if !IsSameKind(subTypeParam, ifcParam) {
				return false
			}
		}
	}
	return true
}

func IsError(typ *types.Type) bool {
	return len(typ.Name.Package) == 0 && typ.Name.Name == "error" && typ.Kind == types.Interface
}

//TODO 内部复杂数据类型还需要再进一步比较
func IsSameKind(typesType *types.Type, reflectType reflect.Type) bool {
	if typesType.Name.Package != reflectType.PkgPath() {
		return false
	}
	result := false
	switch typesType.Kind {
	case types.Builtin:
		switch typesType.Name.Name {
		case "String":
			result = reflectType.Kind() == reflect.String
		case "Int64":
			result = reflectType.Kind() == reflect.Int64
		case "Int32":
			result = reflectType.Kind() == reflect.Int32
		case "Int16":
			result = reflectType.Kind() == reflect.Int16
		case "Int":
			result = reflectType.Kind() == reflect.Int
		case "Uint64":
			result = reflectType.Kind() == reflect.Uint64
		case "Uint32":
			result = reflectType.Kind() == reflect.Uint32
		case "Uint16":
			result = reflectType.Kind() == reflect.Uint16
		case "Uint":
			result = reflectType.Kind() == reflect.Uint
		case "Uintptr":
			result = reflectType.Kind() == reflect.Uintptr
		case "Float64":
			result = reflectType.Kind() == reflect.Float64
		case "Float32":
			result = reflectType.Kind() == reflect.Float32
		case "Float":
			// TODO 这个勉强映射
			result = reflectType.Kind() == reflect.Float64
		case "Bool":
			result = reflectType.Kind() == reflect.Bool
		case "Byte":
			// TODO 这个勉强映射
			result = reflectType.Kind() == reflect.Int8
		}
	case types.Struct:
		result = reflectType.Kind() == reflect.Struct
	case types.Map:
		result = reflectType.Kind() == reflect.Map
	case types.Slice:
		result = reflectType.Kind() == reflect.Slice
	case types.Pointer:
		result = reflectType.Kind() == reflect.Pointer
	case types.Alias:
		//TODO 待映射
	case types.Interface:
		result = reflectType.Kind() == reflect.Interface &&
			len(typesType.Methods) == reflectType.NumMethod()
	case types.Array:
		result = reflectType.Kind() == reflect.Array
	case types.Chan:
		result = reflectType.Kind() == reflect.Chan
	case types.Func:
		result = reflectType.Kind() == reflect.Func
	case types.DeclarationOf:
		//TODO 待映射
	case types.Unknown:
		//TODO 待映射
	case types.Unsupported:
		result = reflectType.Kind() == reflect.Invalid
	case types.Protobuf:
		result = reflectType.Kind() == reflect.Invalid

	}
	return result
}
