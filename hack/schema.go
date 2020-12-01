package hack

import (
	"fmt"
	"reflect"
	"strings"
)

func hackStruct(obj interface{}) {

	v := reflect.ValueOf(obj)
	k := v.Kind()

	if k == reflect.Interface || k == reflect.Ptr {

		v = v.Elem()
		k = v.Kind()
	}

	if k == reflect.Slice || k == reflect.Array {

		for i := 0; i < v.Len(); i++ {

			e := v.Index(i).Interface()
			validateStruct(e)
		}
	} else {
		validateStruct(obj)
	}
}

func validateStruct(obj interface{}) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip ignored and unexported fields in the struct
		if field.Tag.Get("form") == "-" || !val.Field(i).CanInterface() {
			continue
		}

		fieldValue := val.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue) &&
				field.Type.Elem().Kind() == reflect.Struct) {
			validateStruct(fieldValue)
		}

		if true || strings.Index(field.Tag.Get("binding"), "required") > -1 {
			if reflect.DeepEqual(zero, fieldValue) {
				name := field.Name
				if j := field.Tag.Get("json"); j != "" {
					name = j
				} else if f := field.Tag.Get("form"); f != "" {
					name = f
				}
				fmt.Printf("  %s\n", name)
				//errors.Add([]string{name}, RequiredError, "Required")
			}
		}
		//fmt.Printf("--%v\n", fieldValue)
	}
}
