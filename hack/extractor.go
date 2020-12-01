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
	rHandler = reflect.NewAt(rHandler.Type(), unsafe.Pointer(rHandler.UnsafeAddr()))
	q := rHandler.Elem()

	return GetFileLineName(q.Interface())
}

func ExtractRoutesDatarouter(r interface{}) RouteDefinitions {
	var routes []RouteDefinition
	rv := reflect.ValueOf(r) // datarouter.Router (as interface{}) to *datarouter.routerImpl
	rv = reflect.Indirect(rv)
	rRoutes := rv.FieldByName("routeMap")

	key := rRoutes.MapKeys()[0]
	val := rRoutes.MapIndex(key)

	v := reflect.NewAt(val.Type(), unsafe.Pointer(val.InterfaceData()[0]))
	q := v.Elem().Interface()
	fmt.Println(q)

	// structPtr := reflect.NewAt(reflect.StructOf(
	// 	[]reflect.StructField{
	// 		reflect.StructField{
	// 			Name:    "endpoint",
	// 			Type:    reflect.TypeOf("string"),
	// 		},
	// 	},
	// ),
	// 	unsafe.Pointer(val.InterfaceData()[0]))
	// fmt.Println(structPtr.String())
	// structVal := reflect.Indirect(structPtr)
	// fmt.Println(structVal.FieldByName("endpoint"))

	// fmt.Println(val)
	// ptr := unsafe.Pointer(val.InterfaceData()[0])
	// fmt.Println(ptr)
	// v := reflect.NewAt(/* type of struct??? */, unsafe.Pointer(val.InterfaceData()[0]))
	// fmt.Println(reflect.Indirect(v).FieldByName("endpoint"))

	// iter := rRoutes.MapRange()
	// for iter.Next() {
	// 	fmt.Println(iter.Key())
	// 	// rType := iter.Value().Interface()
	// 	rType := iter.Value()
	// 	rValue := reflect.ValueOf(rType)
	// 	rValue = reflect.Indirect(rValue)
	// 	fmt.Println(rValue.FieldByName("endpoint"))
	// 	fmt.Println(reflect.ValueOf(iter.Value()).FieldByName("endpoint"))
	// 	fmt.Println(rType)

	// 	// fmt.Println(rType.Interface())

	// 	// rValue := reflect.ValueOf(rType)

	// 	fmt.Println(reflect.Indirect(rValue))

	// 	r := reflect.Indirect(iter.Value())

	// 	fmt.Println(r.FieldByName("handler"))

	// 	// route := iter.Value()       // datarouter.Route
	// 	// r := reflect.ValueOf(route) // Route -> *datarouter.routeImpl
	// 	// r = reflect.Indirect(r)     // *datarouter.routeImpl -> routeImpl
	// 	// routes = append(routes, collectRouteDatarouter(r))
	// }
	return routes
}

func collectRouteDatarouter(rv reflect.Value) RouteDefinition {
	// fv := reflect.ValueOf(rv).FieldByName("endpoint")
	fv := reflect.ValueOf(rv).FieldByName("handler")
	fmt.Println(fv)
	// for i := 0; i < rv.NumField(); i++ {
	// 	f := rv.Field(i)
	// 	// field := reflect.Indirect(reflect.ValueOf(f))
	// 	// fmt.Println(f.Type().Field(0).Name)
	// 	fmt.Println(reflect.TypeOf(f))
	// }

	rEndpoint := rv.FieldByName("endpoint")
	endpoint := reflect.NewAt(rEndpoint.Type(), unsafe.Pointer(rEndpoint.UnsafeAddr())).Elem().Interface().(string)
	rMethod := rv.FieldByName("method")
	method := reflect.NewAt(rMethod.Type(), unsafe.Pointer(rMethod.UnsafeAddr())).Elem().Interface().(string)
	routeDef := RouteDefinition{
		Method: method,
		Route:  endpoint,
	}

	rHandler := rv.FieldByName("handler")
	file, line, name := getHandlerFuncName(rHandler)
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
