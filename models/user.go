package models


import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
)


type User struct {
	Email string
	Hash []byte
	Salt []byte
}

type UserStore interface {
	Has(email string) bool
	Get(email string) (*User, error)
	Add(u *User) (err error)
	Update(u *User) (err error)
	Remove(u *User) (err error)
}

func Salt() []byte {
	return randBytes(32)
}

func HashPassword(password, salt []byte) (hash []byte, err error) {
	N := 2 << (16 - 1)
	r := 8
	p := 1
	hashLen := 32
	hash, err = scrypt.Key(password, salt, N, r, p, hashLen)
	if err != nil {
		return nil, err
	}
	return hash, err
}

func Login(store UserStore, email, password string) (*User, error) {
	u, err := store.Get(email)
	if err != nil {
		return nil, err
	}
	if !u.VerifyPassword(password) {
		return nil, fmt.Errorf("password invalid")
	}
	return u, nil
}

func Register(store UserStore, email, password string) (*User, error) {
	if store.Has(email) {
		return nil, fmt.Errorf("User with email %v already exists", email)
	}
	u, err := newUser(email, password)
	if err != nil {
		return nil, err
	}
	err = store.Add(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func newUser(email, password string) (*User, error) {
	salt := Salt()
	hash, err := HashPassword([]byte(password), salt)
	if err != nil {
		return nil, err
	}
	u := &User{
		Email: email,
		Hash: hash,
		Salt: salt,
	}
	return u, nil
}

func (u *User) VerifyPassword(attempt string) bool {
	ahash, err := HashPassword([]byte(attempt), u.Salt)
	if err != nil {
		log.Println(err)
		return false
	}
	cmp := subtle.ConstantTimeCompare(u.Hash, ahash)
	return cmp == 1
}

func (u *User) Json() []byte {
	b, err := json.Marshal(&u)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func (u *User) DecodeJson(bytes []byte) error {
	return json.Unmarshal(bytes, &u)
}


