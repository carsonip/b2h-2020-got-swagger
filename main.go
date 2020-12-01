package main

import (
	"encoding/json"
	"fmt"
	"github.com/devfacet/gocmd"
	"io/ioutil"
	"martiniExample/hack"
	"strings"
)

func main() {
	//printRoutes()
	//exportRoutes()
	routes := "/Users/tymon.solecki/dev/pendo-appengine/src/routes.json"
	//matchRoute(routes, "get", "/api/s/5668600916475904/subscription/setting/featureFlags")
	// Init the app
	flags := struct {
		export struct {
			Settings bool `settings:"true" allow-unknown-arg:"true"`
		} `command:"export" description:"Print arguments"`
		Match struct {
			Settings bool `settings:"true" allow-unknown-arg:"true"`
		} `command:"match" description:"Print arguments"`
		Echo      struct {
			Settings bool `settings:"true" allow-unknown-arg:"true"`
		} `command:"echo" description:"Print arguments"`

	}{}

	gocmd.HandleFlag("export", func(cmd *gocmd.Cmd, args []string) error {
		fmt.Printf("exporting to routes.json...")
		exportRoutes()
		return nil
	})
	gocmd.HandleFlag("Echo", func(cmd *gocmd.Cmd, args []string) error {
		fmt.Printf("%s\n", strings.Join(cmd.FlagArgs("Echo")[1:], " "))
		return nil
	})


	gocmd.HandleFlag("Match", func(cmd *gocmd.Cmd, args []string) error {
		matchRoute(routes, cmd.FlagArgs("Match")[1], cmd.FlagArgs("Match")[2])
		return nil
	})

	gocmd.New(gocmd.Options{
		Name:        "basic",
		Version:     "1.0.0",
		Description: "A basic app",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}

func matchRoute(path string, method string, match string) {
	if dat, err := ioutil.ReadFile(path); err != nil {
		fmt.Println("*** Failed to read routes json, maybe you need to generate it first! ***")
	} else {
		routeDefs := hack.RouteDefinitions{}
		json.Unmarshal([]byte(dat), &routeDefs)
		rMatch := routeDefs.MatchPath(method, match)
		lastHandler := rMatch.Handlers[len(rMatch.Handlers)-1]
		fmt.Printf("%v %v\n%v:%v", rMatch.Method, rMatch.Route, lastHandler.Path, lastHandler.LineNo)
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
