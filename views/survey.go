package views

import (
	"fmt"
	"log"
)

import (
)

import (
	"github.com/timtadh/cc-survey/models"
)


func (v *Views) Survey(c *Context) {
	var next string
	err := v.survey.Do(func (s *models.Survey) error {
		cid, _ := s.Next()
		if cid >= 0 {
			next = fmt.Sprintf("/survey/%d", cid)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("was unable to process the request"))
		return
	}
	err = v.tmpl.ExecuteTemplate(c.rw, "survey", map[string]interface{}{
		"email": c.u.Email,
		"clone_count": len(v.clones),
		"next": next,
	})
	if err != nil {
		log.Panic(err)
	}
}

