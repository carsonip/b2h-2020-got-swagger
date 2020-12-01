package hack

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type ContactForm struct {
	Name    string `form:"name" binding:"required"`
	Email   string `form:"email"`
	Message string `form:"message" binding:"required"`
}

func TestGetSchema(t *testing.T) {
	m := martini.NewRouter()
	m.Post("/contact/submit", binding.Bind(ContactForm{}), func(contact ContactForm) string {
		return fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s\n",
			contact.Name, contact.Email, contact.Message)
	})
	rRoutes := getRoutesRv(m)
	rRoute := reflect.Indirect(rRoutes.Index(0))
	rHandlers := getRHandlersFromRRoute(rRoute)
	assert.Equal(t, 2, rHandlers.Len())
	assert.Equal(t, ContactForm{}, getBindStructIfExist(getRHandlerIntFromRHandler(rHandlers.Index(0))))
	assert.Equal(t, nil, getBindStructIfExist(getRHandlerIntFromRHandler(rHandlers.Index(1))))
}
