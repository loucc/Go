// Package main 演示 Go 1.26:new 支持表达式作为初始值。
//
// 1.26 之前:new 只接受类型,得到指向零值的指针。
//
//	p := new(int)   // *int, *p == 0
//	要拿到指向 42 的指针,得写辅助函数或先声明变量再取地址:
//	  x := 42
//	  p := &x
//	或者:
//	  func ptr[T any](v T) *T { return &v }
//	  p := ptr(42)
//
// 1.26 之后:new 接受表达式,直接得到指向该值的指针。
//
//	p := new(42)             // *int,   *p == 42
//	q := new("hello")        // *string, *q == "hello"
//	r := new(User{Name:"a"}) // *User
//
// 最常见的收益:JSON / protobuf 里"可选字段"的指针字段初始化更清爽,
// 不再需要一堆 boolPtr / stringPtr 辅助函数。
package main

import (
	"encoding/json"
	"fmt"
)

type Query struct {
	Page    *int    `json:"page,omitempty"`
	Keyword *string `json:"keyword,omitempty"`
	Active  *bool   `json:"active,omitempty"`
}

func main() {
	// 1.26:直接用表达式初始化
	p := new(42)
	fmt.Printf("p=%T, *p=%d\n", p, *p)

	// 可选字段一行搞定
	q := Query{
		Page:    new(1),
		Keyword: new("golang"),
		Active:  new(true),
	}
	b, _ := json.Marshal(q)
	fmt.Println(string(b))

	// 老写法仍然有效:new(Type) 依旧给零值
	z := new(int)
	fmt.Printf("z=%T, *z=%d (零值)\n", z, *z)
}
