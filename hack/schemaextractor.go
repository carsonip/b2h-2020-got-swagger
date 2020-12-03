package hack

import (
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"net/http"
	"reflect"
	"strings"
)

type MockContext struct {
	bindStruct interface{}
}

func (m *MockContext) Apply(i interface{}) error {
	panic("implement me")
}

func (m *MockContext) Invoke(i interface{}) ([]reflect.Value, error) {
	fun := GetFunctionName(i)
	if !strings.HasPrefix(fun, "github.com/martini-contrib/binding.Form.") {
		return nil, nil
	}
	f, ok := i.(func(ctx martini.Context, r *http.Request))
	if !ok {
		return nil, nil
	}
	r, _ := http.NewRequest(http.MethodGet, "foo", strings.NewReader(""))
	f(m, r)
	return nil, nil
}

func (m *MockContext) Map(i interface{}) inject.TypeMapper {
	if i != nil {
		m.bindStruct = i
	}
	return nil
}

func (m *MockContext) MapTo(i interface{}, i2 interface{}) inject.TypeMapper {
	return nil
}

func (m *MockContext) Set(r reflect.Type, value reflect.Value) inject.TypeMapper {
	panic("implement me")
}

func (m *MockContext) Get(r reflect.Type) reflect.Value {
	return reflect.ValueOf(binding.Errors{})
}

func (m *MockContext) SetParent(injector inject.Injector) {
	panic("implement me")
}

func (m *MockContext) Next() {
	panic("implement me")
}

func (m *MockContext) Written() bool {
	panic("implement me")
}

func getBindStructIfExist(rHandlerInt reflect.Value) interface{} {
	funcName := GetFunctionName(rHandlerInt.Interface())
	if strings.HasPrefix(funcName, "github.com/martini-contrib/binding.Bind.") {
		handlerFunc := rHandlerInt.Elem().Interface().(func(ctx martini.Context, r *http.Request))
		r, _ := http.NewRequest(http.MethodGet, "foo", strings.NewReader(""))
		mc := MockContext{}
		handlerFunc(&mc, r)
		return mc.bindStruct
	}
	return nil
}
