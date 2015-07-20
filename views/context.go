package views

import (
	"log"
	"net/http"
	"runtime/debug"
)

import (
    "github.com/julienschmidt/httprouter"
	"github.com/gorilla/schema"
)

import (
	"github.com/timtadh/cc-survey/models"
)


type Context struct {
	views *Views
	s *models.Session
	u *models.User
	rw http.ResponseWriter
	r *http.Request
	p httprouter.Params
	formDecoder *schema.Decoder
}

func (v *Views) Context(f View) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				log.Println(string(debug.Stack()))
			}
			return
		}()
		c := &Context{
			views: v,
			rw: rw, r: r, p: p,
			formDecoder: schema.NewDecoder(),
		}
		c.Session(v.Log(f))
	}
}

func (v *Views) LoggedOut(f View, to string) View {
	return func(c *Context) {
		if c.u != nil {
			v.sessions.Invalidate(c.s.Key())
			http.Redirect(c.rw, c.r, to, 302)
		} else {
			f(c)
		}
	}
}

func (v *Views) LoggedIn(f View) View {
	return func(c *Context) {
		if c.u == nil {
			c.rw.WriteHeader(401)
			c.rw.Write([]byte("unauthorized"))
		} else {
			f(c)
		}
	}
}

func (v *Views) LoggedInRedirect(f View, to string) View {
	return func(c *Context) {
		if c.u != nil {
			http.Redirect(c.rw, c.r, to, 302)
		} else {
			f(c)
		}
	}
}

func (c *Context) Session(f View) {
	doErr := func(c *Context, err error) {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("error processing request"))
	}
	s, err := models.GetSession(c.views.sessions, c.rw, c.r)
	if err != nil {
		doErr(c, err)
	}
	c.s = s
	if s.User() != "" {
		u, err := c.views.users.Get(s.User())
		if err != nil {
			doErr(c, err)
		}
		c.u = u
	}
	f(c)
}

func (c *Context) SetUser(u *models.User) error {
	c.u = u
	return c.s.SetUser(c.views.sessions, u.Email())
}

