package hack

import (
	"reflect"
	"strings"
)

type FieldType int
const (
	FieldInvalid = iota
	FieldString
	FieldNumber
	FieldInteger
	FieldBoolean
	FieldArray
	FieldObject
)

type Schema struct {
	Name string
	Type FieldType
	ArrayType *Schema
	OmitEmpty bool
	Required bool

	Children []Schema
}

func StructToSchema(obj interface{}) Schema {
	v := reflect.ValueOf(obj)
	k := v.Kind()

	if k == reflect.Interface || k == reflect.Ptr {
		v = v.Elem()
		k = v.Kind()
	}

	s := Schema{}

	if k == reflect.Slice || k == reflect.Array {
		s.Type = FieldArray
		s.ArrayType = getArrayType(v.Type().Elem())
	} else {
		s.Type = FieldObject
		s.Children = validateStruct(obj)
	}
	return s
}

func getArrayType(elem reflect.Type) *Schema {
	var t *Schema
	switch elem.Kind(){
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64:
		t = &Schema{Type: FieldInteger}
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		t = &Schema{Type: FieldNumber}
	case reflect.String:
		t = &Schema{Type: FieldString}
	case reflect.Bool:
		t = &Schema{Type: FieldBoolean}
	case reflect.Slice, reflect.Array:
		t = &Schema{Type: FieldArray}  // TODO: recursive array definition
	}
	return t
}

func validateStruct(obj interface{}) []Schema {
	var s []Schema

	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		f := Schema{}
		field := typ.Field(i)

		// Skip ignored and unexported fields in the struct
		if field.Tag.Get("form") == "-" || field.Tag.Get("json") == "-" || !val.Field(i).CanInterface() {
			continue
		}

		fieldValue := val.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue) &&
				field.Type.Elem().Kind() == reflect.Struct) {
			f.Type = FieldObject
			f.Children = validateStruct(fieldValue)
		} else {
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
				reflect.Uint16, reflect.Uint32, reflect.Uint64:
					f.Type = FieldInteger
			case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
				f.Type = FieldNumber
			case reflect.String:
				f.Type = FieldString
			case reflect.Bool:
				f.Type = FieldBoolean
			case reflect.Slice, reflect.Array:
				f.Type = FieldArray
				f.ArrayType = getArrayType(field.Type.Elem())
			}
		}

		name := field.Name
		if j := field.Tag.Get("json"); j != "" {
			if strings.HasSuffix(j, ",omitempty") {
				j = j[:len(j) - len(",omitempty")]
				f.OmitEmpty = true
			}
			name = j
		} else if f := field.Tag.Get("form"); f != "" {
			name = f
		}
		f.Name = name

		if strings.Index(field.Tag.Get("binding"), "required") > -1 {
			f.Required = true
		}

		s = append(s, f)
	}
	return s
}
