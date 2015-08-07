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
	var a_count int
	var my_count int
	err = v.survey.Do(func (s *models.Survey) error {
		a_count = len(s.Clones) - s.Unanswered.Size()
		my_count = s.CountAnswers(c.u.Email)
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
		"answer_count": a_count,
		"my_answer_count": my_count,
		"next": next,
	})
	if err != nil {
		log.Panic(err)
	}
}

