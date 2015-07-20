package views

import (
	"fmt"
	"log"
	"net/http"
)

import (
	"github.com/gorilla/schema"
)

import (
	"github.com/timtadh/cc-survey/models"
)


func (v *Views) Register(c *Context) {
	err := v.tmpl.ExecuteTemplate(c.rw, "register", map[string]interface{}{
		"target": "/register",
		"csrf": c.s.Csrf("/register"),
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) ErrorRegister(c *Context, r *RegisterForm, errs schema.MultiError) {
	err := v.tmpl.ExecuteTemplate(c.rw, "register", map[string]interface{}{
		"target": "/register",
		"csrf": c.s.Csrf("/register"),
		"errors": errs,
		"form": r,
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) DoRegister(c *Context) {
	r := &RegisterForm{}
	ferr, err := c.Form(r)
	if err != nil {
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed form submitted"))
		return
	} else if len(ferr) > 0 {
		v.ErrorRegister(c, r, ferr)
		return
	}
	_, err = models.Register(v.users, r.Email, r.Password)
	if err != nil {
		ferr["register"] = fmt.Errorf("could not register")
		v.ErrorRegister(c, r, ferr)
		return
	}
	http.Redirect(c.rw, c.r, "/login", 302)
}

type RegisterForm struct {
	Form
	Email string `schema:"email"`
	Password string `schema:"password"`
	ConfirmPassword string `schema:"confirm_password"`
	Consent string `schema:"consent"`
}

func (r *RegisterForm) Validate(c *Context) schema.MultiError {
	errors := r.Form.Validate(c)
	if r.Email == "" {
		errors["email"] = fmt.Errorf("must have an email address")
	}
	if len(r.Password) < 5 {
		errors["password"] = fmt.Errorf("password must be at least 5 characters")
	}
	if r.Password != r.ConfirmPassword {
		errors["confirm_password"] = fmt.Errorf("passwords must match")
	}
	if r.Consent != "iagree" {
		errors["consent"] = fmt.Errorf("You must consent to the survey to register")
	}
	return errors
}
