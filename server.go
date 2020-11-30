package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"reflect"
	"runtime"
	"unsafe"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/foo", foo)
	//if (this is a API tool)
	//{
	//	list all of the routes under m.ROuter.routes
	//}

	v := m.Router

	//fmt.Println("Indirect type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Type()) // prints main.CustomStruct
	//
	//fmt.Println("Indirect value type is:", reflect.Indirect(reflect.ValueOf(v)).Elem().Kind()) // prints struct


	rv := reflect.ValueOf(v)
	fmt.Println(rv.Kind(), rv.Type(), rv)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println(rv.Kind(), rv.Type(), rv)
		rv = rv.Elem()
	}

	fmt.Println(rv.Kind(), rv.Type(), rv)
	rv = reflect.Indirect(rv)
	fmt.Println(rv.Kind(), rv.Type(), rv)
	rv = rv.FieldByName("routes")
	rv = rv.Index(1)
	rv = reflect.Indirect(rv)
	rv = rv.FieldByName("handlers").Index(0)

	rHandler := reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr()))
	//rv = rHandlers.Index(0)
	//handler := rv.Interface().(martini.Handler)
	q := rHandler.Elem()

	//fmt.Println(GetFunctionNameFromReflectValue(q))

	fmt.Println(GetFunctionName(q.Interface()))

	x := q.Interface().(func() string)
	fmt.Println(x)
	//w := reflect.NewAt(q.Type(), unsafe.Pointer(q.UnsafeAddr()))
	//fmt.Println(w)


	//q := rv.Call([]reflect.Value{})
	//fmt.Println(q)

	//x, ok := rv.Interface().(func() string)
	//if !ok {
	//	fmt.Println("yikes")
	//}
	fmt.Println(GetFunctionName(x))


	//
	//route := m.Router.All()[1]
	//
	//v := reflect.ValueOf(route)
	//y := v.FieldByName("handlers")
	//fmt.Println(y.Interface())
	//fmt.Println(GetFunctionName(.handlers[0]))


	//m.Run()
}

func GetFunctionNameFromReflectValue(rv reflect.Value) string {
	return runtime.FuncForPC(rv.Pointer()).Name()
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func foo() string {
	return "foo"
}
