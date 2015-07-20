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


func (v *Views) Logout(c *Context) {
	http.Redirect(c.rw, c.r, "/", 302)
}

func (v *Views) Login(c *Context) {
	err := v.tmpl.ExecuteTemplate(c.rw, "login", map[string]interface{}{
		"target": "/login",
		"csrf": c.s.Csrf("/login"),
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) ErrorLogin(c *Context, l *LoginForm, errs schema.MultiError) {
	err := v.tmpl.ExecuteTemplate(c.rw, "login", map[string]interface{}{
		"target": "/login",
		"csrf": c.s.Csrf("/login"),
		"errors": errs,
		"form": l,
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) DoLogin(c *Context) {
	l := &LoginForm{}
	ferr, err := c.Form(l)
	if err != nil {
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed form submitted"))
		return
	} else if len(ferr) > 0 {
		v.ErrorLogin(c, l, ferr)
		return
	}
	log.Println(l)
	u, err := models.Login(c.views.users, l.Email, l.Password)
	if err != nil {
		log.Println(err)
		ferr["login"] = fmt.Errorf("could not log on to the system. try again?")
		v.ErrorLogin(c, l, ferr)
		return
	} else {
		err := c.SetUser(u)
		if err != nil {
			log.Println(err)
			ferr["login"] = fmt.Errorf("could not log on to the system. try again?")
			v.ErrorLogin(c, l, ferr)
			return
		}
	}
	http.Redirect(c.rw, c.r, "/main", 302)
}

type LoginForm struct {
	Form
	Email string `schema:"email"`
	Password string `schema:"password"`
}

func (l *LoginForm) Validate(c *Context) schema.MultiError {
	errors := l.Form.Validate(c)
	if l.Email == "" {
		errors["email"] = fmt.Errorf("must have an email")
	}
	if len(l.Password) < 1 {
		errors["password"] = fmt.Errorf("must have a password")
	}
	return errors
}
