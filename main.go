package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"martiniExample/hack"
)

func main() {
	//printRoutes()
	//exportRoutes()

	matchRoute("./routes.json", "/users/123")
}

func matchRoute(path string, match string) {
	if dat, err := ioutil.ReadFile(path); err != nil {
		fmt.Println("*** Failed to read routes json, maybe you need to generate it first! ***")
	} else {
		routeDefs := hack.RouteDefinitions{}
		json.Unmarshal([]byte(dat), &routeDefs)
		fmt.Println(routeDefs.MatchPath(match))
	}
}

func printRoutes() {
	m := hack.GetMartini()
	r := m.Router
	routes := hack.ExtractRoutes(r)
	routes.Print()
}

func exportRoutes() {
	m := hack.GetMartini()
	r := m.Router
	routes := hack.ExtractRoutes(r)
	routes.Export()
}
