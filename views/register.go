package views

import (
	"log"
)


func (v *Views) Register(c *Context) {
	err := v.tmpl.ExecuteTemplate(c.rw, "register", nil)
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) DoRegister(c *Context) {
	c.rw.Write([]byte("not-implemented"))
}
