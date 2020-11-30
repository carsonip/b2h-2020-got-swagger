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

	v := m.Router

	rv := reflect.ValueOf(v)
	fmt.Println(rv.Kind(), rv.Type(), rv)

	rv = reflect.Indirect(rv)
	fmt.Println(rv.Kind(), rv.Type(), rv)
	rv = rv.FieldByName("routes")
	rv = rv.Index(1)
	rv = reflect.Indirect(rv)
	rv = rv.FieldByName("handlers").Index(0)

	rHandler := reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr()))
	q := rHandler.Elem()

	fmt.Println(GetFunctionName(q.Interface()))

	//m.Run()
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func foo() string {
	return "foo"
}
