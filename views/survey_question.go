package views

import (
	"fmt"
	"log"
	"strconv"
)

import (
)

import (
	"github.com/timtadh/cc-survey/models"
)

var Questions = []models.Renderable{
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
}

func (v *Views) SurveyQuestion(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed parameter submitted"))
		return
	}
	route := fmt.Sprintf("/survey/%d", cid)
	f := &models.Form{
		Action: route,
		Csrf: c.s.Csrf(route),
		SubmitText: "Submit Answers",
		Questions: Questions,
	}
	err = v.tmpl.ExecuteTemplate(c.rw, "survey_question", map[string]interface{}{
		"cid": cid,
		"clone": v.clones[cid],
		"form": f.HTML(),
	})
	if err != nil {
		log.Panic(err)
	}
}

