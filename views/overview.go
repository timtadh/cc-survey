package views

import (
	"log"
)

import (
)

import (
	"github.com/timtadh/cc-survey/models"
)



func (v *Views) Overview(c *Context) {
	f := &models.Form{
		Action: "/overview",
		Csrf: c.s.Csrf("/overview"),
		SubmitText: "Submit Answers",
		Questions: []models.Renderable{
			&models.MultipleChoice{
				Question: models.Question{
					Name: "howare",
					Question: "How are you today?",
					Required: true,
				},
				Answers: []models.Answer{
					models.Answer{"ok", "I am feeling ok"},
					models.Answer{"bad", "I am feeling bad"},
					models.Answer{"great", "I am feeling GREAT"},
				},
			},
			&models.FreeResponse{
				Question: models.Question{
					Name: "describe",
					Question: "Please decribe your day:",
					Required: false,
				},
				MaxLength: 500,
			},
		},
	}
	err := v.tmpl.ExecuteTemplate(c.rw, "overview", map[string]interface{}{
		"email": c.u.Email,
		"clone_count": len(v.clones),
		"Form": f.HTML(),
	})
	if err != nil {
		log.Panic(err)
	}
}


