// Package main 演示指针。
//
// 学习要点:
//   - Go 有指针,但**没有指针运算**(不能 p++)
//   - &x 取地址,*p 解引用
//   - new(T) 返回 *T 并初始化为零值
//   - Go 1.26:new(expr) 返回指向该表达式值的指针(可选字段/JSON 常用)
//   - 大结构体传参传指针,避免复制
//   - nil 指针解引用会 panic
//   - Go 有 GC,不必手动 free
//
// 运行:go run .
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	// ---- 1. 基本指针 ----
	x := 10
	p := &x // 取地址
	fmt.Printf("x=%d p=%p *p=%d\n", x, p, *p)
	*p = 20 // 通过指针修改
	fmt.Printf("after *p=20, x=%d\n", x)

	// ---- 2. new(T) ----
	pi := new(int) // 分配一个 int 的内存,返回 *int
	*pi = 42
	fmt.Printf("*pi=%d\n", *pi)

	pu := new(User) // 分配一个 User,字段全为零值
	pu.Name = "Alice"
	pu.Age = 30
	fmt.Printf("pu=%+v\n", *pu)

	// ---- 2b. new(表达式)—— Go 1.26 ----
	// 1.26 起,new 的实参可以是表达式,直接得到指向该值的指针。
	// 不再需要 &临时变量 或自己写 func ptr[T](v T) *T { return &v }。
	pj := new(42) // *int, *pj == 42
	fmt.Printf("new(42): *pj=%d\n", *pj)

	// 典型场景:JSON / protobuf 的"可选字段"用指针表示
	type Query struct {
		Page    *int    `json:"page,omitempty"`
		Keyword *string `json:"keyword,omitempty"`
	}
	q := Query{Page: new(1), Keyword: new("golang")}
	fmt.Printf("Query=%+v (*Page=%d, *Keyword=%s)\n", q, *q.Page, *q.Keyword)

	// ---- 3. 值传参 vs 指针传参 ----
	u := User{Name: "Bob", Age: 25}
	renameByValue(u)
	fmt.Printf("值传参后: %+v (未改)\n", u)

	renameByPointer(&u)
	fmt.Printf("指针传参后: %+v\n", u)

	// ---- 4. nil 指针 ----
	var np *int
	fmt.Printf("np == nil: %v\n", np == nil)
	// *np = 100 // 运行时 panic: nil pointer dereference

	// ---- 5. 指针的指针 ----
	y := 100
	pp := &y
	ppp := &pp
	fmt.Printf("**ppp=%d\n", **ppp)

	// ---- 6. 常见规则:何时用指针接收者 ----
	// - 需要修改接收者字段 → 用指针
	// - 结构体较大(> 64 字节) → 用指针避免复制
	// - 保持一致性:同一类型的方法要么全值,要么全指针
}

func renameByValue(u User) {
	u.Name = "Modified"
}

func renameByPointer(u *User) {
	u.Name = "Modified"
}
