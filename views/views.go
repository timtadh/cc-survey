package views

import (
	"io/ioutil"
	"log"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

import (
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/schema"
)

import (
	"github.com/timtadh/cc-survey/clones"
	"github.com/timtadh/cc-survey/models"
	"github.com/timtadh/cc-survey/models/mem"
	"github.com/timtadh/cc-survey/models/file"
)

type Views struct {
	assetPath string
	clonesPath string
	sessions models.SessionStore
	users models.UserStore
	survey models.SurveyStore
	tmpl *template.Template
	clones []*clones.Clone
	decoder *schema.Decoder
}

type View func(*Context)

func signalSelf(s os.Signal) {
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Panic(err)
	}
	err = p.Signal(s)
	if err != nil {
		log.Panic(err)
	}
}

func Routes(assetPath, clonesPath string) http.Handler {
	mux := httprouter.New()
	assetPath = filepath.Clean(assetPath)
	users, err := file.GetUserStore(filepath.Join(assetPath, "data"))
	if err != nil {
		log.Panic(err)
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, os.Kill)
		for s := range sigs {
			users.Close()
			log.Println("Closed Users")
			signal.Stop(sigs)
			signalSelf(s)
			break
		}
	}()
	v := &Views{
		assetPath: assetPath,
		clonesPath: filepath.Clean(clonesPath),
		sessions: mem.NewSessionMapStore("session"),
		users: users,
		decoder: schema.NewDecoder(),
	}
	mux.GET("/", v.Context(v.Index))
	mux.GET("/logout", v.Context(v.LoggedOut(v.Logout, "/")))
	mux.GET("/register", v.Context(v.LoggedInRedirect(v.Register, "/survey")))
	mux.POST("/register", v.Context(v.LoggedInRedirect(v.DoRegister, "/survey")))
	mux.GET("/login", v.Context(v.LoggedInRedirect(v.Login, "/survey")))
	mux.POST("/login", v.Context(v.LoggedInRedirect(v.DoLogin, "/survey")))
	mux.GET("/survey", v.Context(v.LoggedIn(v.Survey)))
	mux.GET("/survey/:clone", v.Context(v.LoggedIn(v.SurveyQuestion)))
	mux.POST("/survey/:clone", v.Context(v.LoggedIn(v.DoSurveyQuestion)))
	
	v.Init()
	return mux
}

func (v *Views) Init() {
	v.loadTemplates()
	v.loadClones()
	var err error
	v.survey, err = file.NewSurveyStore(v.clonesPath, Questions, v.clones)
	if err != nil {
		log.Panic(err)
	}
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
