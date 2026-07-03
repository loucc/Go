// Package main 演示中间件模式:洋葱模型(逐层包装 Handler)。
//
// 学习要点:
//   - HTTP 中间件本质:func(http.Handler) http.Handler
//   - 链式组合:middleware1(middleware2(middleware3(handler)))
//   - 常见中间件:日志、认证、限流、恢复 panic、CORS、请求 ID
//
// 运行:go run .
// 测试:
//
//	curl localhost:8080/hello              # 会被 auth 拒绝
//	curl -H "Authorization: secret" localhost:8080/hello
//	curl -H "Authorization: secret" localhost:8080/boom
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ctxKey string

const reqIDKey ctxKey = "reqID"

// -------- 中间件实现 --------

// 记录访问日志
func loggingMW(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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
}

// 恢复 panic,避免整个服务崩溃
func recoverMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "internal server error", 500)
				fmt.Println("recovered from panic:", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// 校验 Authorization 头
func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "secret" {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 给每个请求生成一个 ID,写入 context
func reqIDMW(next http.Handler) http.Handler {
	var counter int
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		id := fmt.Sprintf("req-%d", counter)
		ctx := context.WithValue(r.Context(), reqIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 组合多个中间件(从右往左包裹)
func chain(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

// -------- Handlers --------

func helloHandler(w http.ResponseWriter, r *http.Request) {
	rid, _ := r.Context().Value(reqIDKey).(string)
	fmt.Fprintf(w, "Hello! reqID=%s\n", rid)
}

func boomHandler(w http.ResponseWriter, r *http.Request) {
	panic("boom!")
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /boom", boomHandler)

	// 中间件顺序:req-id → 日志 → recover → auth → mux
	handler := chain(
		reqIDMW,
		loggingMW(logger),
		recoverMW,
		authMW,
	)(mux)

	logger.Info("start", "addr", ":8080")
	_ = http.ListenAndServe(":8080", handler)
}
