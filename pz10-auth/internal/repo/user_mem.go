package repo

import (
	"errors"
	"example.com/pz10-auth/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  []byte
}

// Реализуем методы интерфейса UserInfo
func (u UserRecord) GetID() int64     { return u.ID }
func (u UserRecord) GetEmail() string { return u.Email }
func (u UserRecord) GetRole() string  { return u.Role }

type UserMem struct{ users map[string]UserRecord }

func NewUserMem() *UserMem {
	hash := func(s string) []byte {
		h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
		return h
	}
	return &UserMem{users: map[string]UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
	}}
}

var ErrNotFound = errors.New("user not found")
var ErrBadCreds = errors.New("bad credentials")

func (r *UserMem) ByEmail(email string) (UserRecord, error) {
	u, ok := r.users[email]
	if !ok {
		return UserRecord{}, ErrNotFound
	}
	return u, nil
}

func (r *UserMem) CheckPassword(email, pass string) (core.UserInfo, error) {
	u, err := r.ByEmail(email)
	if err != nil {
		return nil, ErrNotFound
	}
	if bcrypt.CompareHashAndPassword(u.Hash, []byte(pass)) != nil {
		return nil, ErrBadCreds
	}
	return u, nil
}
