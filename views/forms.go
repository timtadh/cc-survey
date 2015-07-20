package views

import (
	"fmt"
)

import (
	"github.com/gorilla/schema"
)

type Validatable interface {
	Validate(*Context) schema.MultiError
}

type Form struct {
	Csrf string `schema:"csrf"`
}

func (f *Form) Validate(c *Context) schema.MultiError {
	errors := make(schema.MultiError)
	if !c.s.ValidCsrf(c.r.URL.Path, f.Csrf) {
		errors["csrf"] = fmt.Errorf("invalid csrf token")
	}
	return errors
}

func (c *Context) Form(v Validatable) (schema.MultiError, error) {
	err := c.r.ParseForm()
	if err != nil {
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed form submitted"))
		return nil, fmt.Errorf("could not process form")
	}
	err = c.views.decoder.Decode(v, c.r.PostForm)
	if err != nil {
		return err.(schema.MultiError), nil
	}
	return v.Validate(c), nil
}

