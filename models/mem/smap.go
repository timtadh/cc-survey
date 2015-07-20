package mem

import (
	"fmt"
	"sync"
)

import (
)

import (
	"github.com/timtadh/cc-survey/models"
)


type SessionMapStore struct {
	lock sync.Mutex
	name string
	store map[uint64]*models.Session
}

func NewSessionMapStore(name string) *SessionMapStore {
	return &SessionMapStore{
		name: name,
		store: make(map[uint64]*models.Session),
	}
}

func (m *SessionMapStore) Name() string {
	return m.name
}

func (m *SessionMapStore) Get(key uint64) (*models.Session, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if s, has := m.store[key]; has {
		return s.Copy(), nil
	}
	return nil, fmt.Errorf("Session not in store")
}

func (m *SessionMapStore) Invalidate(key uint64) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.store, key)
	return nil
}

func (m *SessionMapStore) Update(s *models.Session) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if s == nil {
		return fmt.Errorf("passed in a nil session")
	}
	m.store[s.Key()] = s.Copy()
	return nil
}

