package models

import (
	"bytes"
	"html/template"
	"log"
)


type Question struct {
	Name string
	Text string
	Required bool
}

type MultipleChoice struct {
	Question
	Answers []Answer
}

type Answer struct {
	Name string
	Text string
}
var aTmpl = template.Must(template.New("answer").Parse(`<input type="radio" id="{{.group}}" name="{{.group}}" value="{{.Name}}"> {{.text}}`))

type FreeResponse struct {
	Question
	MaxLength int
}

type Templatable interface {
	Template() string
}

type Renderable interface {
	HTML() template.HTML
}

func (q *FreeResponse) HTML() template.HTML {
	buf := new(bytes.Buffer)
	err := q.Template().Execute(buf, q)
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

func (q *FreeResponse) Template() *template.Template {
	t := `<label for="{{.Name}}">
	<div class="question{{if .Required}} required{{end}}">{{.Text}}</div>
	<textarea id="{{.Name}}" name="{{.Name}}" maxlength={{.MaxLength}}>
	</textarea>
</label>`
	return template.Must(template.New(q.Name).Parse(t))
}

func (q *MultipleChoice) HTML() template.HTML {
	buf := new(bytes.Buffer)
	err := q.Template().Execute(buf, q)
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

func (q *MultipleChoice) Template() *template.Template {
	t := `<label for="{{.Name}}">
<div class="question{{if .Required}} required{{end}}">
{{.Text}}
</div>
{{with $q := .}} {{range $answer := .Answers}}<div class="answer">
	{{$answer.HTML $q.Name}}
</div>{{end}}{{end}}
</label>`
	return template.Must(template.New(q.Name).Parse(t))
}

func (a *Answer) HTML(group string) template.HTML {
	buf := new(bytes.Buffer)
	err := aTmpl.Execute(buf, map[string]interface{}{
		"group": group,
		"name": a.Name,
		"text": a.Text,
	})
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

