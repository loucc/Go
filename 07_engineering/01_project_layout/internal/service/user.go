// Package service 演示业务逻辑层(不依赖 HTTP/DB 细节)。
package service

import (
	"errors"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID   int64
	Name string
}

// UserService 是纯业务对象。数据来源可通过接口注入(此处内存实现)。
type UserService struct {
	mu    sync.RWMutex
	users map[int64]User
	next  int64
}

func New() *UserService {
	return &UserService{
		users: map[int64]User{},
		next:  1,
	}
}

func (s *UserService) Create(name string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: s.next, Name: name}
	s.next++
	s.users[u.ID] = u
	return u
}

func (s *UserService) Get(id int64) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (s *UserService) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}
