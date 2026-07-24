// Package main 演示 fmt 格式化动词。
//
// 学习要点:
//   - %v 通用值;%+v 显示字段名;%#v 显示 Go 语法表示
//   - %T 类型;%p 指针
//   - %d/%b/%o/%x:整数进制
//   - %f/%e/%g:浮点
//   - %s/%q:字符串;%c 字符
//   - %w:错误包装(仅在 fmt.Errorf 中)
//   - Sprintf 返回字符串,Fprintf 写到 Writer
//
// 运行:go run .
package main

import (
	"fmt"
	"os"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{"Alice", 30}

	// ---- 打印结构体 ----
	fmt.Printf("%%v  : %v\n", u)  // {Alice 30}
	fmt.Printf("%%+v : %+v\n", u) // {Name:Alice Age:30}
	fmt.Printf("%%#v : %#v\n", u) // main.User{Name:"Alice", Age:30}
	fmt.Printf("%%T  : %T\n", u)  // main.User

	// ---- 整数 ----
	n := 255
	fmt.Printf("dec=%d bin=%b oct=%o hex=%x HEX=%X\n", n, n, n, n, n)

	// 宽度、填充、对齐
	fmt.Printf("[%5d]\n", 42)  // [   42] 右对齐
	fmt.Printf("[%-5d]\n", 42) // [42   ] 左对齐
	fmt.Printf("[%05d]\n", 42) // [00042] 补零

	// ---- 浮点 ----
	f := 3.14159
	fmt.Printf("f=%f\n", f)   // 3.141590
	fmt.Printf("f=%.2f\n", f) // 3.14
	fmt.Printf("f=%e\n", f)   // 3.141590e+00
	fmt.Printf("f=%g\n", f)   // 3.14159

	// ---- 字符串 ----
	s := "hello\tworld"
	fmt.Printf("%%s : %s\n", s) // 原样
	fmt.Printf("%%q : %q\n", s) // 带引号并转义

	// ---- rune ----
	r := 'A'
	fmt.Printf("%%c=%c  %%U=%U\n", r, r) // A  U+0041

	// ---- 指针 ----
	fmt.Printf("%%p : %p\n", &u)

	// ---- Sprintf / Fprintf ----
	msg := fmt.Sprintf("user=%s age=%d", u.Name, u.Age)
	fmt.Println("Sprintf:", msg)
	fmt.Fprintln(os.Stderr, "写到 stderr:", msg)

	// ---- 错误包装 %w ----
	base := fmt.Errorf("db closed")
	wrapped := fmt.Errorf("query failed: %w", base)
	fmt.Println("wrapped:", wrapped)
}
