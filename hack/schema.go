package hack

import (
	"fmt"
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

// https://swagger.io/docs/specification/data-models/data-types/
func (ft FieldType) String() string {
	ftToStrMap := map[FieldType]string {
		FieldInvalid: "INVALID",
		FieldString: "string",
		FieldNumber: "number",
		FieldInteger: "integer",
		FieldBoolean: "boolean",
		FieldArray: "array",
		FieldObject: "object",
	}
	return ftToStrMap[ft]
}

type Schema struct {
	Name      string
	Type      FieldType
	ArrayType *Schema
	OmitEmpty bool
	Required  bool
	Recursive bool

	Children []Schema
}

func (s Schema) Print() {
	s.printIndent(0)
}

func (s Schema) printIndent(indent int) {
	indentStr := strings.Repeat(" ", indent)
	fmt.Printf("%s%s %s\n", indentStr, s.Name, s.Type)
	if s.ArrayType != nil {
		fmt.Printf("%sArrayType:\n", indentStr)
		s.ArrayType.printIndent(indent + 2)
	} else if len(s.Children) > 0 {
		fmt.Printf("%sChildren:\n", indentStr)
		for _, c := range s.Children {
			c.printIndent(indent + 2)
		}
	}
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
		s.ArrayType = getArrayType(v.Type().Elem(), nil)
	} else {
		s.Type = FieldObject
		s.Children = traverseStruct(obj)
	}
	return s
}

func getArrayType(elem reflect.Type, parentType reflect.Type) *Schema {
	var t *Schema
	kind := elem.Kind()
	if kind == reflect.Ptr {
		elem = elem.Elem()
		kind = elem.Kind()
	}
	zero := reflect.Zero(elem).Interface()
	switch kind {
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
		t = &Schema{Type: FieldArray, ArrayType: getArrayType(elem.Elem(), elem)}
	case reflect.Struct:
		if elem == parentType {
			t = &Schema{Type: FieldObject, Recursive: true}
		} else {
			t = &Schema{Type: FieldObject, Children: traverseStruct(zero)}
		}
	case reflect.Map:
		t = &Schema{Type: FieldObject}
	}
	return t
}

func traverseStruct(obj interface{}) []Schema {
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

		kind := field.Type.Kind()
		if kind == reflect.Struct {
			f.Type = FieldObject
			f.Children = traverseStruct(fieldValue)
		} else if kind == reflect.Ptr && field.Type.Elem() == val.Type() {
			f.Type = FieldObject
			f.Recursive = true
		} else if kind == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			f.Type = FieldObject
			f.Children = traverseStruct(reflect.Zero(field.Type.Elem()).Interface())
		} else {
			if kind == reflect.Ptr {
				kind = field.Type.Elem().Kind()
			}
			switch kind {
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
				if field.Type.Elem().Kind() == reflect.Uint8 {
					// []byte should be a string
					f.Type = FieldString
				} else {
					f.Type = FieldArray
					f.ArrayType = getArrayType(field.Type.Elem(), typ)
				}
			case reflect.Map:
				f.Type = FieldObject
			}
		}

		name := field.Name
		if j := field.Tag.Get("json"); j != "" {
			if strings.HasSuffix(j, ",omitempty") {
				j = j[:len(j)-len(",omitempty")]
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
