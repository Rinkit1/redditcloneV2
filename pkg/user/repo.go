package user

import (
	"errors"
	"strconv"
	"sync"
)

var (
	ErrNoUser           = errors.New("no user found")
	ErrBadPass          = errors.New("invalid password")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type UserMemoryRepository struct {
	data   map[string]*User
	lastID int
	mu     *sync.RWMutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data:   map[string]*User{},
		lastID: 0,
		mu:     &sync.RWMutex{},
	}
}
func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	repo.mu.Lock()
	u, ok := repo.data[login]
	repo.mu.Unlock()
	if !ok {
		return nil, ErrNoUser
	}
	if u.password != pass {
		return nil, ErrBadPass
	}
	return u, nil
}
func (repo *UserMemoryRepository) AddUser(login, pass string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, ok := repo.data[login]
	if ok {
		return ErrUserAlreadyExist
	}
	repo.lastID++
	repo.data[login] = &User{
		ID:       strconv.Itoa(repo.lastID),
		Login:    login,
		password: pass,
	}
	return nil
}
