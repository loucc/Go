// Package main 演示 Go 1.22 增强的 net/http 路由。
//
// 新特性:
//   - 方法前缀:  "GET /users"
//   - 路径变量:  "/users/{id}"
//   - 通配符:    "/files/{path...}"
//   - r.PathValue("id") 提取
//
// 运行:go run .
// 测试:
//
//	curl localhost:8080/users/42
//	curl -X POST localhost:8080/users/42
//	curl localhost:8080/files/a/b/c.txt
package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET user %s\n", r.PathValue("id"))
	})
	mux.HandleFunc("POST /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "POST user %s\n", r.PathValue("id"))
	})

	// 通配符 {name...} 匹配剩余所有段
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "file path: %s\n", r.PathValue("path"))
	})

	fmt.Println("listening :8080")
	_ = http.ListenAndServe(":8080", mux)
}
