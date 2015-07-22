package views

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

import (
	"github.com/gorilla/schema"
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
	if err != nil || cid < 0 || cid >= len(v.clones) {
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

func (v *Views) ErrorSurveyQuestion(c *Context, cid int, f *models.Form, a *models.SurveyAnswer, errs schema.MultiError) {
	err := v.tmpl.ExecuteTemplate(c.rw, "survey_question", map[string]interface{}{
		"cid": cid,
		"clone": v.clones[cid],
		"form": f.HTML(),
		"errors": errs,
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) DoSurveyQuestion(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil || cid < 0 || cid >= len(v.clones) {
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
	answer, ferr, err := f.Decode(c.u, cid, c.r)
	if err != nil {
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed form submitted"))
		return
	} else if len(ferr) > 0 {
		v.ErrorSurveyQuestion(c, cid, f, answer, ferr)
		return
	}
	var next string
	err = v.survey.Do(func (s *models.Survey) error {
		s.Answer(answer)
		cid, _ := s.Next()
		if cid >= 0 {
			next = fmt.Sprintf("/survey/%d", cid)
		} else {
			next = "/survey"
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("there was an error processing your answer"))
		return
	}
	http.Redirect(c.rw, c.r, next, 302)
}


