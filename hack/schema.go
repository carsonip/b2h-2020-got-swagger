package hack

import (
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"net/http"
	"reflect"
	"strings"
)

type MockContext struct {}

func (m MockContext) Apply(i interface{}) error {
	panic("implement me")
}

func (m MockContext) Invoke(i interface{}) ([]reflect.Value, error) {
	fun := GetFunctionName(i)
	if !strings.HasPrefix(fun, "github.com/martini-contrib/binding.Form.") {
		return nil, nil
	}
	f, ok := i.(func(ctx martini.Context, r *http.Request))
	if !ok {
		return nil, nil
	}
	r, _ := http.NewRequest(http.MethodGet, "foo", strings.NewReader(""))
	f(MockContext{}, r)
	return nil, nil
}

func (m MockContext) Map(i interface{}) inject.TypeMapper {
	if i != nil {
		hackStruct(i)
	}
	return nil
}

func (m MockContext) MapTo(i interface{}, i2 interface{}) inject.TypeMapper {
	return nil
}

func (m MockContext) Set(r reflect.Type, value reflect.Value) inject.TypeMapper {
	panic("implement me")
}

func (m MockContext) Get(r reflect.Type) reflect.Value {
	return reflect.ValueOf(binding.Errors{})
}

func (m MockContext) SetParent(injector inject.Injector) {
	panic("implement me")
}

func (m MockContext) Next() {
	panic("implement me")
}

func (m MockContext) Written() bool {
	panic("implement me")
}

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

func getBindStructIfExist(rHandlerInt reflect.Value) {
	funcName := GetFunctionName(rHandlerInt.Interface())
	if strings.HasPrefix(funcName, "github.com/martini-contrib/binding.Bind.") {
		handlerFunc := rHandlerInt.Elem().Interface().(func(ctx martini.Context, r *http.Request))
		r, _ := http.NewRequest(http.MethodGet, "foo", strings.NewReader(""))
		fmt.Println("fields:")
		handlerFunc(MockContext{}, r)
	}
}

