package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"martiniExample/hack"
)

func main() {
	printRoutes()
	//exportRoutes()

	//matchRoute("./routes.json", "get", "/api/s/5668600916475904/subscription/setting/featureFlags")
}

func matchRoute(path string, method string, match string) {
	if dat, err := ioutil.ReadFile(path); err != nil {
		fmt.Println("*** Failed to read routes json, maybe you need to generate it first! ***")
	} else {
		routeDefs := hack.RouteDefinitions{}
		json.Unmarshal([]byte(dat), &routeDefs)
		rMatch := routeDefs.MatchPath(method, match)
		rJSON, err := json.Marshal(rMatch)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("%s\n", string(rJSON))
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
