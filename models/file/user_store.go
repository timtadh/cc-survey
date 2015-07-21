package file

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sync"
)

import (
	"github.com/timtadh/fs2/fmap"
	"github.com/timtadh/fs2/bptree"
)

import (
	"github.com/timtadh/cc-survey/models"
)


type UserFileStore struct {
	lock sync.RWMutex
	path string
	bf *fmap.BlockFile
	users *bptree.BpTree
}

func GetUserStore(dir string) (*UserFileStore, error) {
	var bf *fmap.BlockFile
	var users *bptree.BpTree
	path := filepath.Join(dir, "users.bptree")
	err := createOrOpen(path,
		func(path string) (err error) {
			bf, err = fmap.CreateBlockFile(path)
			if err != nil {
				return err
			}
			users, err = bptree.New(bf, -1, -1)
			return err
		},
		func(path string) (err error) {
			bf, err = fmap.OpenBlockFile(path)
			if err != nil {
				return err
			}
			users, err = bptree.Open(bf)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	s := &UserFileStore{
		path: path,
		bf: bf,
		users: users,
	}
	return s, bf.Sync()
}

func (s *UserFileStore) Close() error {
	return s.bf.Close()
}

func (s *UserFileStore) Has(email string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.has(email)
}

func (s *UserFileStore) has(email string) bool {
	has, err := s.users.Has([]byte(email))
	if err != nil {
		log.Panic(err)
	}
	return has
}

func (s *UserFileStore) Get(email string) (u *models.User, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	err = s.users.DoFind([]byte(email), func(_, bytes []byte) error {
		u = &models.User{}
		return json.Unmarshal(bytes, &u)
	})
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("User not found")
	}
	return u, nil
}

func (s *UserFileStore) Add(u *models.User) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	b, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	return s.add(u.Email, b)
}

func (s *UserFileStore) add(email string, user []byte) (err error) {
	if s.has(email) {
		return fmt.Errorf("store already has user")
	}
	return s.users.Add([]byte(email), user)
}

func (s *UserFileStore) Remove(u *models.User) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.remove(u.Email)
}

func (s *UserFileStore) remove(email string) (err error) {
	return s.users.Remove([]byte(email), func(_ []byte) bool { return true })
}

func (s *UserFileStore) Update(u *models.User) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	err = s.remove(u.Email)
	if err != nil {
		return err
	}
	b, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	return s.add(u.Email, b)
}


