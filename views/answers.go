package views

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
)

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

import (
	"github.com/timtadh/cc-survey/models"
)


func (v *Views) Answers(c *Context) {
	var err error
	var a_count int
	err = v.survey.Do(func (s *models.Survey) error {
		a_count = len(s.Clones) - s.Unanswered.Size()
		return nil
	})
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("was unable to process the request"))
		return
	}
	err = v.tmpl.ExecuteTemplate(c.rw, "answers", map[string]interface{}{
		"clone_count": len(v.clones),
		"answer_count": a_count,
		"next": "/answers/0",
	})
	if err != nil {
		log.Panic(err)
	}
}


func (v *Views) Answer(c *Context) {
	var err error
	aid, err := strconv.Atoi(c.p.ByName("answer"))
	if err != nil || aid < 0 {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("invalid answer id")
		}
		c.rw.WriteHeader(400)
		c.rw.Write([]byte("malformed parameter submitted"))
		return
	}
	var next string
	var prev string
	if aid > 0 {
		prev = fmt.Sprintf("/answers/%d", aid - 1)
	}
	var answer *models.SurveyAnswer
	err = v.survey.Do(func (s *models.Survey) error {
		if aid >= len(s.Answers) {
			return fmt.Errorf("invalid answer id")
		}
		if aid + 1 < len(s.Answers) {
			next = fmt.Sprintf("/answers/%d", aid + 1)
		}
		answer = s.Answers[aid]
		return nil
	})
	if err != nil {
		log.Println(err)
		c.rw.WriteHeader(500)
		c.rw.Write([]byte("was unable to process the request"))
		return
	}
	err = v.tmpl.ExecuteTemplate(c.rw, "answer", map[string]interface{}{
		"answer": answer,
		"questions": Questions,
		"cid": answer.CloneID,
		"clone": v.clones[answer.CloneID],
		"next": next,
		"prev": prev,
		"mark": func(answer string) template.HTML {
			unsafe := blackfriday.MarkdownCommon([]byte(answer))
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			return template.HTML(string(html))
		},
	})
	if err != nil {
		log.Panic(err)
	}
}
