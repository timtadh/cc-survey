package views

import (
	"log"
)

import (
)

import (
	"github.com/timtadh/cc-survey/models"
)



func (v *Views) Main(c *Context) {
	q := &models.MultipleChoice{
		Question: models.Question{
			Name: "howare",
			Text: "How are you today?",
			Required: true,
		},
		Answers: []models.Answer{
			models.Answer{"ok", "I am feeling ok"},
			models.Answer{"bad", "I am feeling bad"},
			models.Answer{"great", "I am feeling GREAT"},
		},
	}
	err := v.tmpl.ExecuteTemplate(c.rw, "main", map[string]interface{}{
		"email": c.u.Email,
		"clone_count": len(v.clones),
		"Question": q.HTML(),
	})
	if err != nil {
		log.Panic(err)
	}
}


