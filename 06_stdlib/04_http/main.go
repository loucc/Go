// Package main 演示 net/http:一个最简单的 REST 服务。
//
// Go 1.22+ 的路由器支持方法+通配符:
//
//	mux.HandleFunc("GET /users/{id}", ...)
//
// 学习要点:
//   - http.Handler / http.HandlerFunc
//   - http.ServeMux 路由
//   - r.PathValue("id") 提取路径参数
//   - 优雅关闭 http.Server.Shutdown
//   - 中间件模式:http.Handler 是可组合的
//   - 安全实践:超时配置、请求体大小限制
//
// 运行:go run .
// 测试:
//
//	curl -X GET  http://localhost:8080/users/1
//	curl -X POST http://localhost:8080/users -d '{"name":"Alice"}'
//	curl -X GET  http://localhost:8080/health
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	mu    sync.RWMutex
	users = map[int64]User{
		1: {ID: 1, Name: "Alice"},
		2: {ID: 2, Name: "Bob"},
	}
	nextID int64 = 3
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// Go 1.22+ 方法+路径通配
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("POST /users", createUser)

	// 用中间件包一层
	handler := loggingMiddleware(logger, mux)

	// 安全实践:配置超时,防止 Slowloris 等慢速攻击
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,  // 读取请求头的超时
		ReadTimeout:       10 * time.Second, // 读取整个请求体的超时
		WriteTimeout:      15 * time.Second, // 写入响应的超时
		IdleTimeout:       60 * time.Second, // keep-alive 连接空闲超时
		MaxHeaderBytes:    1 << 20,          // 请求头最大 1 MB
	}

	// 启动
	go func() {
		logger.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// 等待信号,优雅关闭
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown err", "err", err)
	}
	logger.Info("bye")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // 1.22+
	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	mu.RLock()
	u, ok := users[id]
	mu.RUnlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// 安全实践:限制请求体大小,防止恶意大文件攻击
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad json or body too large", http.StatusBadRequest)
		return
	}
	mu.Lock()
	u.ID = nextID
	nextID++
	users[u.ID] = u
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// 中间件:记录每次请求
func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("access",
			"method", r.Method,
			"path", r.URL.Path,
			"dur", time.Since(start),
		)
	})
}
