package hack

import "github.com/go-martini/martini"

func GetMartini() *martini.ClassicMartini {
	m := martini.Classic()
	m.Get("/.*", func() string {
		return "Hello world!"
	})

	m.Group("/api", func(publicApiRouter martini.Router) {
		publicApiRouter.Get("/get", func() string {
			return "hello world"
		})

		m.Group("/dashboard", func(publicApiRouter martini.Router) {
			m.Group("/share", func(publicApiRouter martini.Router) {
				publicApiRouter.Get("/list", func() string {
					return "list handler"
				})
			})
		})
	})
	return m
}

func foo() string {
	return "foo"
}


