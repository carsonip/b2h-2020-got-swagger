package main

import (
	"./hack"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	m := hack.GetMartini()
	r := m.Router
	routes := hack.ExtractRoutes(r)
	printer(routes)
	//m.Run()
}

func printer(routes []hack.RouteDefinition) {
	pwd, _ := os.Getwd()
	for _, r := range routes {
		fmt.Println(r.Method, r.Route)
		for _, h := range r.Handlers {
			relPath, _ := filepath.Rel(pwd, h.Path)
			relFuncName, _ := filepath.Rel(pwd, h.FuncName)
			fmt.Printf("    %v:%v %v\n", relPath, h.LineNo, relFuncName)
		}
	}
}
