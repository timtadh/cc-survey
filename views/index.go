package views

import (
	"html/template"
	"net/http"
)

import (
	"github.com/julienschmidt/httprouter"
)

var index = template.Must(template.New("/").Parse(`
<html>
<body>
	<h1>Hello</h1>
</body>
</html>
`))

func (v *Views) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	index.Execute(rw, nil)
}

