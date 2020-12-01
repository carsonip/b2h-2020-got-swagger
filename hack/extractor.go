package hack

import (
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"unsafe"
	"encoding/json"
)

type routeHandler struct {
	Path     string `json:"path"`
	LineNo   int `json:"lineNo"`
	FuncName string `json:"funcName"`
}

type RouteDefinition struct {
	Method string `json:"method"`
	Route string `json:"route"`
	Handlers []routeHandler `json:"handlers"`
}

type RouteDefinitions []RouteDefinition

func (routes RouteDefinitions) Print() {
	pwd, _ := os.Getwd()
	for _, r := range routes {
		fmt.Println(r.Method, r.Route)
		for _, h := range r.Handlers {
			relPath, _ := filepath.Rel(pwd, h.Path)
			fmt.Printf("    %v:%v %v\n", relPath, h.LineNo, h.FuncName)
		}
	}
}

func (routes RouteDefinitions) Export() {
	jsonRoutes, _ := json.Marshal(routes)
	if err := ioutil.WriteFile("./routes.json", jsonRoutes, 0644); err != nil {
		fmt.Println("*** Failed to write file ***")
	}
	//fmt.Println(string(jsonRoutes))
}

func ExtractRoutes(r martini.Router) RouteDefinitions {
	var routes []RouteDefinition

	rv := reflect.ValueOf(r)  // Router interface to *router
	rv = reflect.Indirect(rv) // *router to router
	rRoutes := rv.FieldByName("routes")
	for i := 0; i < rRoutes.Len(); i++ {
		routes = append(routes, collectRoute(rRoutes.Index(i)))
	}
	return routes
}

func collectRoute(rv reflect.Value) RouteDefinition {
	rRoute := reflect.Indirect(rv) // *route to route

	rPattern := rRoute.FieldByName("pattern")
	pattern := reflect.NewAt(rPattern.Type(), unsafe.Pointer(rPattern.UnsafeAddr())).Elem().Interface().(string)
	rMethod := rRoute.FieldByName("method")
	method := reflect.NewAt(rMethod.Type(), unsafe.Pointer(rMethod.UnsafeAddr())).Elem().Interface().(string)
	rHandlers := rRoute.FieldByName("handlers")


	return newRoute(method, pattern, rHandlers)
}

var routeReg1 = regexp.MustCompile(`:[^/#?()\.\\]+`)
var routeReg2 = regexp.MustCompile(`\*\*`)

func newRoute(method string, pattern string, rHandlers reflect.Value) RouteDefinition {
	routeDef := RouteDefinition{
		Method: method,
		Route:  pattern,
	}

	for i := 0; i < rHandlers.Len(); i++ {
		file, line, name := getHandlerFuncName(rHandlers.Index(i))
		routeDef.Handlers = append(routeDef.Handlers, routeHandler{
			Path:     file,
			LineNo:   line,
			FuncName: name,
		})
	}

	pattern = routeReg1.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	var index int
	pattern = routeReg2.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	return routeDef
}

func getHandlerFuncName(rHandler reflect.Value) (string, int, string) {
	rHandlerPtr := reflect.NewAt(rHandler.Type(), unsafe.Pointer(rHandler.UnsafeAddr()))
	rHandlerInt := rHandlerPtr.Elem()
	file, line, funcName := GetFileLineName(rHandlerInt.Interface())
	getBindStructIfExist(rHandlerInt)
	return file, line, funcName
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetFileLineName(i interface{}) (string, int, string) {
	pc := reflect.ValueOf(i).Pointer()
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	funcName := f.Name()
	if funcName[0:1] == "_" {
		funcName = funcName[1:]
	}
	return file, line, funcName
}
