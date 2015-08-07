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
			Name: "duplicated-code",
			Question: "(1) Are the snippets above duplicated code?",
			Required: true,
		},
		Answers: []models.Answer{
			models.Answer{"yes", "Yes"},
			models.Answer{"no", "No"},
		},
	},
	&models.MultipleChoice{
		Question: models.Question{
			Name: "why-not-duplicated",
			Question: "(2) If you answered *No* to question 1, why do you consider the above code not duplicated?",
			Required: false,
		},
		Answers: []models.Answer{
			models.Answer{"one-line", "It is only one line of code"},
			models.Answer{"for-loop", "It seems to be just highlighting a for loop or iterator protocol"},
			models.Answer{"string-builder", "It is just calls to StringBuilder"},
			models.Answer{"standard-lib", "It is only calls to the standard library (either Java or Android)"},
			models.Answer{"reflection", "It just a call to a reflected method"},
			models.Answer{"other", "Other"},
		},
	},
	&models.FreeResponse{
		Question: models.Question{
			Name: "why-other-not-duplicated",
			Question: "(3) If you answered *Other* to question 2, explain why:",
			Required: false,
		},
		MaxLength: 1000,
	},
	&models.MultipleChoice{
		Question: models.Question{
			Name: "action-to-take",
			Question: "(4) If you answered *Yes* to question 1, would you:",
			Required: false,
		},
		Answers: []models.Answer{
			models.Answer{"create-story", "Create a story to refactor this code"},
			models.Answer{"comment", "Add a comment to consider refactoring on next change"},
			models.Answer{"comment-no-refactor", "Add a note about duplicate code even if it can't be refactored"},
			models.Answer{"ignore", "Ignore it"},
			models.Answer{"action", "Take some other action"},
		},
	},
	&models.FreeResponse{
		Question: models.Question{
			Name: "why-ignore-clones",
			Question: "(5) If you answered *Ignore It* to question 4, explain why:",
			Required: false,
		},
		MaxLength: 1000,
	},
	&models.FreeResponse{
		Question: models.Question{
			Name: "why-take-action",
			Question: "(6) If you answered *Take some other action* to question 4, explain why:",
			Required: false,
		},
		MaxLength: 1000,
	},
	&models.MultipleChoice{
		Question: models.Question{
			Name: "pattern-characteristics",
			Question: "(7) If you answered *Yes* to question 1, do you consider this pattern to be",
			Required: false,
		},
		Answers: []models.Answer{
			models.Answer{"example", "A good example of how to do something"},
			models.Answer{"only-way", "The only way to do something"},
			models.Answer{"best-way", "The best example of how to do something"},
			models.Answer{"fine-way", "The neither a good or bad example of how to do something"},
			models.Answer{"bad-way", "A bad example of how to do something"},
			models.Answer{"incorrect", "An incorrect example of how to do something"},
			models.Answer{"not-example", "None of the above"},
		},
	},
	&models.FreeResponse{
		Question: models.Question{
			Name: "other-thoughts",
			Question: "(8) If you answered *Yes* to question 1, do you have any other thoughts on this code?",
			Required: false,
		},
		MaxLength: 1000,
	},
}

func (v *Views) SurveyQuestion(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil || cid < 0 || cid >= len(v.clones) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid clone id")
		}
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
		"form": f.HTML(map[string]error{}, map[string]string{}),
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) ErrorSurveyQuestion(c *Context, cid int, f *models.Form, a *models.SurveyAnswer, errs schema.MultiError, answers map[string]string) {
	err := v.tmpl.ExecuteTemplate(c.rw, "survey_question", map[string]interface{}{
		"cid": cid,
		"clone": v.clones[cid],
		"form": f.HTML(errs, answers),
	})
	if err != nil {
		log.Panic(err)
	}
}

func (v *Views) DoSurveyQuestion(c *Context) {
	cid, err := strconv.Atoi(c.p.ByName("clone"))
	if err != nil || cid < 0 || cid >= len(v.clones) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid clone id")
		}
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
	answer, ferr, answers, err := f.Decode(c.s, c.u, v.clones[cid], cid, c.r)
	if err != nil {
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed form submitted"))
		return
	} else if len(ferr) > 0 {
		v.ErrorSurveyQuestion(c, cid, f, answer, ferr, answers)
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


