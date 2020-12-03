package main

import (
	"encoding/json"
	"fmt"
	"github.com/devfacet/gocmd"
	"io/ioutil"
	"martiniExample/hack"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

func main() {
	//printRoutes()
	//exportRoutes()
	routes := "./routes.json"

	// Init the app
	flags := struct {
		export struct {
			Settings bool `settings:"true" allow-unknown-arg:"true"`
		} `command:"export" description:"export the routers from server.go to Routes.json"`
		Match struct {
			Routes   string `short:"r" long:"routes" description:"Path to your Routes.json file"`
			Method   string `short:"m" long:"method" description:"request Method, e.g. POST"`
			Path     string `short:"p", long:"path", description:"request Path, e.g./api/some/Path"`
			Headers  string `short:"h" long:"headers" description:"request Headers copied from the browser in the format: :Method: <Method> :Path: <Path>, e.g. :Method: POST :Path: /api/aggregation/s/5630785994358784"`
			Settings bool   `settings:"true" allow-unknown-arg:"true"`
		} `command:"match" description:"match the request. Example requests:\n./martiniExample match -m get -p /api/s/5668600916475904/subscription/setting/featureFlag, \n./martiniExample match -h \":Method: POST :Path: /api/aggregation/s/5630785994358784\"\n --- "`
		Init struct {
			Settings bool   `settings:"true" allow-unknown-arg:"true"`
		} `command:"init" description:"Generates routes.json"`
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

	gocmd.HandleFlag("Init", func(cmd *gocmd.Cmd, args []string) error {
		fmt.Printf("%s\n", strings.Join(cmd.FlagArgs("Init")[1:], " "))

		buildRoutesJson()

		return nil
	})


	gocmd.HandleFlag("Match", func(cmd *gocmd.Cmd, args []string) error {
		method := ""
		path := ""
		if flags.Match.Headers == "" && (flags.Match.Method == "" || flags.Match.Path == "") {
			panic("Supply either Headers flag or Method and Path flags!")
		}
		if flags.Match.Headers != "" {
			methodFromCurlRegex, _ := regexp.Compile(":Method: *([\\w/]+)")
			pathFromCurlRegex, _ := regexp.Compile(":Path: *([\\w/]+)")
			methodMatch := methodFromCurlRegex.FindAllStringSubmatch(cmd.FlagArgs("Match")[1], -1)
			pathMatch := pathFromCurlRegex.FindAllStringSubmatch(cmd.FlagArgs("Match")[1], -1)
			method = methodMatch[0][1]
			path = pathMatch[0][1]
		} else {
			method = flags.Match.Method
			path = flags.Match.Path
		}
		if flags.Match.Routes != ""{
			routes = flags.Match.Routes
		}

		if _, err := os.Stat(routes); os.IsNotExist(err) {
			buildRoutesJson()
		}

		matchRoute(routes, method, path)

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
		fmt.Println(err.Error())
		fmt.Println("*** Failed to read routes.json, maybe you need to generate it first! ***")
	} else {
		routeDefs := hack.RouteDefinitions{}
		json.Unmarshal([]byte(dat), &routeDefs)
		rMatch := routeDefs.MatchPath(method, match)
		lastHandler := rMatch.Handlers[len(rMatch.Handlers)-1]
		fmt.Printf("\n  \033[0;32m%v\033[0;0m %v\n\n Is handled by:\n  > %v:%v\n", rMatch.Method, rMatch.Route, lastHandler.Path, lastHandler.LineNo)
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


func buildRoutesJson() {
		usr, _ := user.Current()

		dat, err := ioutil.ReadFile(usr.HomeDir + "/dev/pendo-appengine/src/pendo.io/server/server.go")
		if err != nil {
			panic(err.Error())
		}

		newFile := strings.Replace(string(dat), "/* HACK HERE */", "x := hack.ExtractRoutes(router)\n\tx.Export()\n\tpanic(\"foo\")", 1)
		if newFile == string(dat) {
			panic("Nothing changed - did you add the magic comment?")
		}

		newFile = strings.Replace(newFile, "import (", "import (\n\t\"github.com/carsonip/b2h/hack\"", 1)
		ioutil.WriteFile(usr.HomeDir + "/dev/pendo-appengine/src/pendo.io/server/server.go", []byte(newFile), 0644)

		c := exec.Command("go", "build", "-o", "appEngineHack", ".")
		c.Dir = usr.HomeDir + "/dev/pendo-appengine/src/appengine"
		c.Run()


		c2 := exec.Command("./appEngineHack")
		c2.Dir = usr.HomeDir + "/dev/pendo-appengine/src/appengine"
		c2.Run()

		c3 := exec.Command( "mv", usr.HomeDir + "/dev/pendo-appengine/src/appengine/routes.json", ".")
		c3.Run()

		c4 := exec.Command("rm", "appEngineHack")
		c4.Dir = usr.HomeDir + "/dev/pendo-appengine/src/appengine"
		c4.Run()

		ioutil.WriteFile(usr.HomeDir + "/dev/pendo-appengine/src/pendo.io/server/server.go", []byte(dat), 0644)
}
