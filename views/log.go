package views

import (
	"log"
	"net/http"
	"time"
)

import (
    "github.com/julienschmidt/httprouter"
)

type loggingRW struct {
	rw http.ResponseWriter
	total int
}

func (l *loggingRW) Header() http.Header {
	return l.rw.Header()
}

func (l *loggingRW) Write(bytes []byte) (int, error) {
	c, err := l.rw.Write(bytes)
	l.total += c
	return c, err
}

func (l *loggingRW) WriteHeader(code int) {
	l.rw.WriteHeader(code)
}

func (c *Context) Log(f httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		lrw := &loggingRW{rw: rw}
		s := time.Now()
		f(lrw, r, p)
		e := time.Now()
		log.Printf("%v %v (%v) %v (%d) %v", r.RemoteAddr, r.URL, r.ContentLength, c.s.Key(), lrw.total, e.Sub(s))
	}
}

