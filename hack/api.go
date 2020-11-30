package hack

import "github.com/go-martini/martini"

func GetMartini() *martini.ClassicMartini {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/foo", foo)
	return m
}

func foo() string {
	return "foo"
}
