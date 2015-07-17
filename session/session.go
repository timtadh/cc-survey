package session

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
    "github.com/julienschmidt/httprouter"
)

type Session struct {
	key uint64
	csrf uint64
	addr string
	usrAgent string
	created time.Time
	accessed time.Time
}

type Store interface {
	Name() string
	Get(key uint64) (*Session, error)
	Update(*Session) (error)
	Invalidate(key uint64) (error)
	Session(f func(*Session)httprouter.Handle) httprouter.Handle
}

func randBytes(length int) []byte {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		log.Fatal(err)
	} else {
		slice := make([]byte, length)
		if _, err := urandom.Read(slice); err != nil {
			log.Fatal(err)
		}
		urandom.Close()
		return slice
	}
	panic("unreachable")
}

func randUint64() uint64 {
	b := randBytes(8)
	return binary.LittleEndian.Uint64(b)
}

func userAgent(r *http.Request) string {
	if agent, has := r.Header["User-Agent"]; has {
		return strings.Join(agent, "; ")
	}
	return "None"
}

func ip(r *http.Request) string {
	return strings.SplitN(r.RemoteAddr, ":", 2)[0]
}

func key(name string, r *http.Request) (uint64, error) {
	c, err := r.Cookie(name)
	if err == nil {
		n, err := strconv.ParseUint(c.Value, 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, fmt.Errorf("Failed to extract session key")
}

func Get(store Store, rw http.ResponseWriter, r *http.Request) (s *Session, err error) {
	name := store.Name()
	k, err := key(name, r)
	if err != nil {
		s = newSession(r)
	} else {
		s, err = store.Get(k)
		if err != nil {
			log.Println(err)
			s = newSession(r)
		} else {
			err := s.update(name, r)
			if err != nil {
				log.Println(err)
				s = newSession(r)
			}
		}
	}
	err = store.Update(s)
	if err != nil {
		return nil, err
	}
	s.write(name, rw, r)
	return s, nil
}


func newSession(r *http.Request) *Session {
	return &Session{
		key: randUint64(),
		csrf: randUint64(),
		addr: ip(r),
		usrAgent: userAgent(r),
		created: time.Now().UTC(),
		accessed: time.Now().UTC(),
	}
}

func (s *Session) Copy() *Session {
	return &Session{
		key: s.key,
		csrf: s.csrf,
		addr: s.addr,
		usrAgent: s.usrAgent,
		created: s.created,
		accessed: s.accessed,
	}
}

func (s *Session) Key() uint64 {
	return s.key
}

func (s *Session) Invalidate(store Store, rw http.ResponseWriter) error {
	delete(rw.Header(), store.Name())
	return store.Invalidate(s.key)
}

func (s *Session) valid(name string, r *http.Request) bool {
	k, err := key(name, r)
	if err != nil {
		return false
	}
	ua := userAgent(r)
	addr := ip(r)
	return ua == s.usrAgent && addr == s.addr && k == s.key
}

func (s *Session) update(name string, r *http.Request) error {
	if s.valid(name, r) {
		s.accessed = time.Now().UTC()
		return nil
	}
	return fmt.Errorf("session was invalid")
}

func (s *Session) write(name string, rw http.ResponseWriter, r *http.Request) {
	v := strconv.FormatUint(s.key, 10)
	secure := r.URL.Scheme == "https" || r.TLS != nil
	http.SetCookie(rw, &http.Cookie{
		Name: name,
		Value: v,
		Path: "/",
		Secure: secure,
		HttpOnly: true,
	})
}

