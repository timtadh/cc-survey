package views

import (
	"log"
)

import (
)

import (
)



func (v *Views) Main(c *Context) {
	err := v.tmpl.ExecuteTemplate(c.rw, "main", map[string]interface{}{
		"email": c.u.Email(),
		"clone_count": len(v.clones),
	})
	if err != nil {
		log.Panic(err)
	}
}


