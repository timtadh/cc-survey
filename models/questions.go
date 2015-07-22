package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

import (
	"github.com/gorilla/schema"
)


type Question struct {
	Name string
	Question string
	Required bool
}

type MultipleChoice struct {
	Question
	Answers []Answer
}

type Answer struct {
	Value string
	Answer string
}

type FreeResponse struct {
	Question
	MaxLength int
}

type Form struct {
	Action string
	Csrf string
	SubmitText string
	Questions []Renderable
}

type Renderable interface {
	HTML() template.HTML
}

var freeTmpl = template.Must(template.New("freeResponse").Parse(
`<label>
	<div class="question{{if .Required}} required{{end}}">
		{{.Question.Question}}
	</div>
	<div class="answer">
		<textarea id="{{.Name}}" name="{{.Name}}" maxlength={{.MaxLength}}></textarea>
	</div>
</label>`))

var multiTmpl = template.Must(template.New("multipleChoice").Parse(
`<label>
	<div class="question{{if .Required}} required{{end}}">
		{{.Question.Question}}
	</div>{{with $q := .}}{{range $a := .Answers}}
	<div class="answer">
		<input type="radio" name="{{$q.Name}}" value="{{$a.Value}}"/>
		{{$a.Answer}}
	</div>{{end}}{{end}}
</label>`))

var formTmpl = template.Must(template.New("form").Parse(
`<form action="{{.Action}}" method="post">{{range $q := .Questions}}
{{$q.HTML}}{{end}}
<input type="hidden" name="csrf" value="{{.Csrf}}"/>
<div class="submit"><input type="submit" value="{{.SubmitText}}"/></div>
</form>`))


func (f *Form) Decode(u *User, cid int, r *http.Request) (*SurveyAnswer, schema.MultiError, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, nil, err
	}
	answer := &SurveyAnswer{
		UserEmail: u.Email,
		CloneID: cid,
		Responses: make([]Response, 0, len(f.Questions)),
	}
	errors := make(schema.MultiError)
	form := r.PostForm
	for qid, r := range f.Questions {
		switch q := r.(type) {
		case *MultipleChoice:
			if value, has := form[q.Name]; !has && q.Required {
				errors[q.Name] = fmt.Errorf("This is a required question")
				answer.Responses = append(answer.Responses, Response{
					QuestionID: qid,
					Answer: -1,
					Text: "Not Answered",
				})
			} else if !has {
				answer.Responses = append(answer.Responses, Response{
					QuestionID: qid,
					Answer: -1,
					Text: "Not Answered",
				})
			} else {
				aid, err := q.AnswerNumber(value[0])
				if err != nil {
					errors[q.Name] = err
					answer.Responses = append(answer.Responses, Response{
						QuestionID: qid,
						Answer: -2,
						Text: "Bad Answer",
					})
				} else {
					answer.Responses = append(answer.Responses, Response{
						QuestionID: qid,
						Answer: aid,
						Text: value[0],
					})
				}
			}
		case *FreeResponse:
			value, has := form[q.Name]
			has = has && value[0] != ""
			if !has && q.Required {
				errors[q.Name] = fmt.Errorf("This is a required question")
				answer.Responses = append(answer.Responses, Response{
					QuestionID: qid,
					Answer: -1,
					Text: "Not Answered",
				})
			} else if !has {
				answer.Responses = append(answer.Responses, Response{
					QuestionID: qid,
					Answer: -1,
					Text: "Not Answered",
				})
			} else {
				answer.Responses = append(answer.Responses, Response{
					QuestionID: qid,
					Answer: -3,
					Text: value[0],
				})
			}
		default:
			log.Panic(fmt.Errorf("unexpected question type"))
		}
	}
	return answer, errors, nil
}

func (f *Form) HTML() template.HTML {
	return HTML(formTmpl, f)
}

func (q *FreeResponse) HTML() template.HTML {
	return HTML(freeTmpl, q)
}

func (q *MultipleChoice) HTML() template.HTML {
	return HTML(multiTmpl, q)
}

func (q *MultipleChoice) AnswerNumber(key string) (int, error) {
	for aid, a := range q.Answers {
		if key == a.Value {
			return aid, nil
		}
	}
	return -1, fmt.Errorf("Not a valid answer")
}

func HTML(t *template.Template, data interface{}) template.HTML {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

