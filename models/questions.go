package models

import (
	"bytes"
	"html/template"
	"log"
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


func (f *Form) HTML() template.HTML {
	return HTML(formTmpl, f)
}

func (q *FreeResponse) HTML() template.HTML {
	return HTML(freeTmpl, q)
}

func (q *MultipleChoice) HTML() template.HTML {
	return HTML(multiTmpl, q)
}

func HTML(t *template.Template, data interface{}) template.HTML {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

