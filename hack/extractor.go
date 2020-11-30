package hack

import (
	"fmt"
	"github.com/go-martini/martini"
	"reflect"
	"runtime"
	"unsafe"
)

func ExtractRoutes(r martini.Router) {
	rv := reflect.ValueOf(r)  // Router interface to *router
	rv = reflect.Indirect(rv)  // *router to router
	rRoutes := rv.FieldByName("routes")
	for i := 0; i < rRoutes.Len(); i++ {
		printRoute(rRoutes.Index(i))
	}

}

func printRoute(rv reflect.Value) {
	rRoute := reflect.Indirect(rv)  // *route to route

	rPattern := rRoute.FieldByName("pattern")
	pattern := reflect.NewAt(rPattern.Type(), unsafe.Pointer(rPattern.UnsafeAddr())).Elem().Interface().(string)
	rMethod := rRoute.FieldByName("method")
	method := reflect.NewAt(rMethod.Type(), unsafe.Pointer(rMethod.UnsafeAddr())).Elem().Interface().(string)
	rHandlers := rRoute.FieldByName("handlers")

	fmt.Printf("%s %s\n", pattern, method)
	for i := 0; i < rHandlers.Len(); i++ {
		file, line, name := getHandlerFuncName(rHandlers.Index(i))
		fmt.Printf("  %s:%d %s\n", file, line, name)
	}
}

func getHandlerFuncName(rHandler reflect.Value) (string, int, string) {
	rHandler = reflect.NewAt(rHandler.Type(), unsafe.Pointer(rHandler.UnsafeAddr()))
	q := rHandler.Elem()

	return GetFileLineName(q.Interface())
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetFileLineName(i interface{}) (string, int, string) {
	pc := reflect.ValueOf(i).Pointer()
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	return file, line, f.Name()
}
