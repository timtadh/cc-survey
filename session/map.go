package session

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

import (
    "github.com/julienschmidt/httprouter"
)


type MapStore struct {
	lock sync.Mutex
	name string
	store map[uint64]*Session
}

func NewMapStore(name string) *MapStore {
	return &MapStore{
		name: name,
		store: make(map[uint64]*Session),
	}
}

func (m *MapStore) Session(f func(*Session) httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s, err := Get(m, rw, r)
		if err == nil {
			f(s)(rw, r, p)
		} else {
			log.Println(err)
			rw.WriteHeader(500)
			rw.Write([]byte("error processing request"))
		}
	}
}

func (m *MapStore) Name() string {
	return m.name
}

func (m *MapStore) Get(key uint64) (*Session, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if s, has := m.store[key]; has {
		return s.Copy(), nil
	}
	return nil, fmt.Errorf("Session not in store")
}

func (m *MapStore) Invalidate(key uint64) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.store, key)
	return nil
}

func (m *MapStore) Update(s *Session) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if s == nil {
		return fmt.Errorf("passed in a nil session")
	}
	m.store[s.key] = s.Copy()
	return nil
}

