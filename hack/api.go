package hack

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type ContactForm struct {
	Name    string `form:"name" binding:"required"`
	Email   string `form:"email"`
	Message string `form:"message" binding:"required"`
}

func GetMartini() *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/contact/submit", binding.Bind(ContactForm{}), func(contact ContactForm) string {
		return fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s\n",
			contact.Name, contact.Email, contact.Message)
	})
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


