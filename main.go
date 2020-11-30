package main

import "martiniExample/hack"

func main() {
	m := hack.GetMartini()
	r := m.Router
	hack.ExtractRoutes(r)
	//m.Run()
}