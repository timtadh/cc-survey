package views

import (
	"log"
)


func (v *Views) Index(c *Context) {
	log.Println(c.u)
	err := v.tmpl.ExecuteTemplate(c.rw, "index", nil)
	if err != nil {
		log.Panic(err)
	}
}

