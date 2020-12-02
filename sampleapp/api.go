package sampleapp

import (
	"github.com/go-martini/martini"
)

func GetMartini() *martini.ClassicMartini {
	m := martini.Classic()
	m.Get("/foo", foo)
	m.Get("/users/:id", foo)

	m.Group("/api", func(publicApiRouter martini.Router) {
		publicApiRouter.Get("/get", func() string {
			return "hello world"
		})
	})
	return m
}

func foo() string {
	return "foo"
}


