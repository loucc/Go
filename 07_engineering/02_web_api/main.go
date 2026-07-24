// Package main 演示 REST API 全流程(单文件,便于阅读)。
//
// 端点:
//
//	GET  /users        列出所有用户
//	POST /users        创建用户
//	GET  /users/{id}   获取
//	PUT  /users/{id}   更新
//	DELETE /users/{id} 删除
//
// 运行:go run .
// 测试:
//
//	curl localhost:8080/users
//	curl -X POST localhost:8080/users -d '{"name":"Alice"}'
//	curl localhost:8080/users/1
//	curl -X PUT localhost:8080/users/1 -d '{"name":"Alice2"}'
//	curl -X DELETE localhost:8080/users/1
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Store struct {
	mu   sync.RWMutex
	data map[int64]User
	next int64
}

var ErrNotFound = errors.New("not found")

func NewStore() *Store {
	return &Store{data: map[int64]User{}, next: 1}
}

func (s *Store) Create(name string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: s.next, Name: name}
	s.next++
	s.data[u.ID] = u
	return u
}

func (s *Store) Get(id int64) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.data[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (s *Store) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.data))
	for _, u := range s.data {
		out = append(out, u)
	}
	return out
}

func (s *Store) Update(id int64, name string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.data[id]
	if !ok {
		return User{}, ErrNotFound
	}
	u.Name = name
	s.data[id] = u
	return u, nil
}

func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := NewStore()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.List())
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Name string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			writeError(w, 400, "name required")
			return
		}
		u := store.Create(body.Name)
		writeJSON(w, http.StatusCreated, u)
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, 400, "invalid id")
			return
		}
		u, err := store.Get(id)
		if err != nil {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, u)
	})
	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, 400, "invalid id")
			return
		}
		var body struct{ Name string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		u, err := store.Update(id, body.Name)
		if err != nil {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, u)
	})
	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, 400, "invalid id")
			return
		}
		if err := store.Delete(id); err != nil {
			writeError(w, 404, "not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	logger.Info("starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("serve failed", "err", err)
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": fmt.Sprintf("%d: %s", code, msg)})
}
