package views

import (
	"io/ioutil"
	"log"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

import (
    "github.com/julienschmidt/httprouter"
)

import (
	"github.com/timtadh/cc-survey/clones"
	"github.com/timtadh/cc-survey/session"
)

type Views struct {
	assetPath string
	clonesPath string
	sessions session.Store
	tmpl *template.Template
	clones []*clones.Clone
}

type Context struct {
	v *Views
	s *session.Session
}


func Routes(assetPath, clonesPath string) http.Handler {
	mux := httprouter.New()
	v := &Views{
		assetPath: filepath.Clean(assetPath),
		clonesPath: filepath.Clean(clonesPath),
		sessions: session.NewMapStore("session"),
	}
	mux.GET("/", v.Log(v.sessions.Session(func(s *session.Session) httprouter.Handle { return v.Context(s).Index })))
	v.Init()
	return mux
}

func (v *Views) Init() {
	v.loadTemplates()
	v.loadClones()
}

func (v *Views) Context(s *session.Session) *Context {
	return &Context{v, s}
}

func (v *Views) loadClones() {
	c, err := clones.LoadAll(v.clonesPath)
	if err != nil {
		log.Panic(err)
	}
	v.clones = c
	log.Println("loaded clones", len(v.clones))
}

func (v *Views) loadTemplates() {
	s, err := os.Stat(v.assetPath)
	if os.IsNotExist(err) {
		log.Fatalf("Could not load assets from %v. Path does not exist.", v.assetPath)
	} else if err != nil {
		log.Panic(err)
	}
	v.tmpl = template.New("!")
	if s.IsDir() {
		v.loadTemplatesFromDir("", filepath.Join(v.assetPath, "templates"), v.tmpl)
	} else {
		log.Fatalf("Could not load assets from %v. Unknown file type", v.assetPath)
	}
}

func (v *Views) loadTemplatesFromDir(ctx, path string, t *template.Template) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}
	for _, info := range dir {
		c := filepath.Join(ctx, info.Name())
		p := filepath.Join(path, info.Name())
		if info.IsDir() {
			v.loadTemplatesFromDir(c, p, t)
		} else {
			v.loadTemplateFile(ctx, p, t)
		}
	}
}

func (v *Views) loadTemplateFile(ctx, path string, t *template.Template) {
	name := filepath.Base(path)
	if strings.HasPrefix(name, ".") {
		return
	}
	ext := filepath.Ext(name)
	if ext != "" {
		name = strings.TrimSuffix(name, ext)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	v.loadTemplate(filepath.Join(ctx, name), string(content), t)
}

func (v *Views) loadTemplate(name, content string, t *template.Template) {
	log.Println("loaded template", name)
	_, err := t.New(name).Parse(content)
	if err != nil {
		log.Panic(err)
	}
}
