package main

import (
	"martiniExample/hack"
)

func main() {
	m := hack.GetMartini()
	r := m.Router
	routes := hack.ExtractRoutes(r)
	routes.Print()
	//m.Run()
}
