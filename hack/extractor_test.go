package hack

import (
	"github.com/go-martini/martini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/", func() string {
		return "Hello world!"
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "GET")
	assert.Equal(t, defs[0].Route, "/")
}

func TestPost(t *testing.T) {
	m := martini.NewRouter()
	m.Post("/", func() string {
		return "Hello world!"
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "POST")
	assert.Equal(t, defs[0].Route, "/")
}

func TestAny(t *testing.T) {
	m := martini.NewRouter()
	m.Any("/", func() string {
		return "Hello world!"
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "*")
	assert.Equal(t, defs[0].Route, "/")
}

func TestGroup(t *testing.T) {
	m := martini.NewRouter()
	m.Group("/api", func(r martini.Router) {
		r.Get("/hello", func() string {
			return "hello world"
		})
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "GET")
	assert.Equal(t, defs[0].Route, "/api/hello")
}

func TestMultipleRoutes(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/foo", func() string {
		return "foo"
	})
	m.Get("/bar", func() string {
		return "bar"
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 2)
	assert.Equal(t, defs[0].Method, "GET")
	assert.Equal(t, defs[0].Route, "/foo")
	assert.Equal(t, defs[1].Method, "GET")
	assert.Equal(t, defs[1].Route, "/bar")
}

func TestNamedParams(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/:id", func(params martini.Params) string {
		return params["id"]
	})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "GET")
	assert.Equal(t, defs[0].Route, "/:id")
}

func TestGlobRoute(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/**", func() {})
	defs := ExtractRoutes(m)
	assert.Len(t, defs, 1)
	assert.Equal(t, defs[0].Method, "GET")
	assert.Equal(t, defs[0].Route, "/**")
}

func TestHandler(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/", func() {})
	defs := ExtractRoutes(m)
	assert.Len(t, defs[0].Handlers, 1)
}

func TestMultiHandler(t *testing.T) {
	m := martini.NewRouter()
	m.Get("/", func() {}, func() {})
	defs := ExtractRoutes(m)
	assert.Len(t, defs[0].Handlers, 2)
}

func TestGroupHandler(t *testing.T) {
	m := martini.NewRouter()
	m.Group("/api", func(r martini.Router) {
		r.Get("/hello", func() {})
	}, func() {})
	defs := ExtractRoutes(m)
	assert.Len(t, defs[0].Handlers, 2)
}
