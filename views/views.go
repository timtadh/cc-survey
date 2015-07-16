package views

import (
	"net/http"
)

import (
    "github.com/julienschmidt/httprouter"
)

type Views struct {
}

func Routes() http.Handler {
	mux := httprouter.New()
	v := &Views{}
	mux.GET("/", v.Index)
	return mux
}
