package views

import (
	"log"
	"net/http"
)

import (
	"github.com/julienschmidt/httprouter"
)


func (v *Views) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := v.tmpl.ExecuteTemplate(rw, "index", nil)
	if err != nil {
		log.Panic(err)
	}
}

