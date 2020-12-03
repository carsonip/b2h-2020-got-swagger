package hack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"unsafe"

	"github.com/go-martini/martini"
)

type routeHandler struct {
	Path     string `json:"path"`
	LineNo   int    `json:"lineNo"`
	FuncName string `json:"funcName"`
}

type RouteDefinition struct {
	Method   string         `json:"method"`
	Route    string         `json:"route"`
	Handlers []routeHandler `json:"handlers"`
	Schema   Schema         `json:"-"`
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

	rRoutes := getRoutesRv(r)
	for i := 0; i < rRoutes.Len(); i++ {
		routes = append(routes, collectRoute(rRoutes.Index(i)))
	}
	return routes
}

func getRoutesRv(r martini.Router) reflect.Value {
	rv := reflect.ValueOf(r)  // Router interface to *router
	rv = reflect.Indirect(rv) // *router to router
	return rv.FieldByName("routes")
}

func collectRoute(rRoutePtr reflect.Value) RouteDefinition {
	rRoute := reflect.Indirect(rRoutePtr) // *route to route

	pattern := getPatternFromRRoute(rRoute)
	method := getMethodFromRRoute(rRoute)

	rHandlers := getRHandlersFromRRoute(rRoute)

	return newRoute(method, pattern, rHandlers)
}

var routeReg1 = regexp.MustCompile(`:[^/#?()\.\\]+`)
var routeReg2 = regexp.MustCompile(`\*\*`)

func newRoute(method string, pattern string, rHandlers reflect.Value) RouteDefinition {
	routeDef := RouteDefinition{
		Method: method,
		Route:  pattern,
	}

	var bindStruct interface{}
	for i := 0; i < rHandlers.Len(); i++ {
		file, line, name, handlerBindStruct := getHandlerMetadata(rHandlers.Index(i))
		if bindStruct == nil && handlerBindStruct != nil {
			bindStruct = handlerBindStruct
		}
		routeDef.Handlers = append(routeDef.Handlers, routeHandler{
			Path:     file,
			LineNo:   line,
			FuncName: name,
		})
	}
	if bindStruct != nil {
		routeDef.Schema = StructToSchema(bindStruct)
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

func getPatternFromRRoute(rRoute reflect.Value) string {
	rPattern := rRoute.FieldByName("pattern")
	return reflect.NewAt(rPattern.Type(), unsafe.Pointer(rPattern.UnsafeAddr())).Elem().Interface().(string)
}

func getMethodFromRRoute(rRoute reflect.Value) string {
	rMethod := rRoute.FieldByName("method")
	return reflect.NewAt(rMethod.Type(), unsafe.Pointer(rMethod.UnsafeAddr())).Elem().Interface().(string)
}

func getRHandlersFromRRoute(rRoute reflect.Value) reflect.Value {
	return rRoute.FieldByName("handlers")
}

func getHandlerMetadata(rHandler reflect.Value) (string, int, string, interface{}) {
	rHandlerInt := getRHandlerIntFromRHandler(rHandler)
	file, line, funcName := GetFileLineName(rHandlerInt.Interface())
	bindStruct := getBindStructIfExist(rHandlerInt)
	return file, line, funcName, bindStruct
}

func getRHandlerIntFromRHandler(rHandler reflect.Value) reflect.Value {
	rHandlerPtr := reflect.NewAt(rHandler.Type(), unsafe.Pointer(rHandler.UnsafeAddr()))
	return rHandlerPtr.Elem()
}

func ExtractRoutesDatarouter(r interface{}) RouteDefinitions {
	var routes []RouteDefinition
	rv := reflect.ValueOf(r)
	rv = reflect.Indirect(rv)
	rRoutes := rv.FieldByName("routeMap")

	iter := rRoutes.MapRange()
	for iter.Next() {
		routeVal := iter.Value().Elem()
		routeVal = reflect.Indirect(routeVal)
		routes = append(routes, collectRouteDatarouter(routeVal))
	}

	return routes
}

func collectRouteDatarouter(rv reflect.Value) RouteDefinition {
	rEndpoint := rv.FieldByName("endpoint")
	endpoint := reflect.NewAt(rEndpoint.Type(), unsafe.Pointer(rEndpoint.UnsafeAddr())).Elem().Interface().(string)
	rMethod := rv.FieldByName("method")
	method := reflect.NewAt(rMethod.Type(), unsafe.Pointer(rMethod.UnsafeAddr())).Elem().Interface().(string)
	routeDef := RouteDefinition{
		Method: method,
		Route:  endpoint,
	}

	rHandler := rv.FieldByName("handler")
	file, line, name := GetFileLineName(rHandler)
	routeDef.Handlers = append(routeDef.Handlers, routeHandler{
		Path:     file,
		LineNo:   line,
		FuncName: name,
	})

	return routeDef
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
