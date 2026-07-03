// Package main 演示结构体。
//
// 学习要点:
//   - struct 是零值可用的,不需要构造函数
//   - 字段标签(tag)配合反射,常用于 JSON/DB 序列化
//   - 匿名字段(嵌入)= 组合
//   - 结构体是值类型,但可以通过指针共享
//   - 空结构体 struct{} 占 0 字节,用作信号或 set
//
// 运行:go run .
package main

import (
	"encoding/json"
	"fmt"
)

// 1. 基本结构体
type Point struct {
	X, Y int
}

// 2. 带 tag 的结构体(常用于序列化)
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`               // 序列化时忽略
	Email    string `json:"email,omitempty"` // 空值时不输出
}

// 3. 匿名字段(嵌入)——组合
type Address struct {
	City, Street string
}

type Employee struct {
	Name    string
	Salary  float64
	Address // 匿名嵌入
}

func main() {
	// ---- 1. 多种初始化方式 ----
	p1 := Point{1, 2}        // 位置初始化(不推荐,顺序敏感)
	p2 := Point{X: 3, Y: 4}  // 字段名初始化(推荐)
	p3 := Point{}            // 零值
	p4 := &Point{X: 5, Y: 6} // 指针
	fmt.Printf("p1=%v p2=%v p3=%v p4=%v\n", p1, p2, p3, *p4)

	// ---- 2. 匿名结构体 ----
	person := struct {
		Name string
		Age  int
	}{Name: "Alice", Age: 30}
	fmt.Printf("匿名结构体: %+v\n", person)

	// ---- 3. 嵌入的字段提升 ----
	e := Employee{
		Name:    "Bob",
		Salary:  10000,
		Address: Address{City: "Beijing", Street: "长安街"},
	}
	// 可以直接访问嵌入字段
	fmt.Printf("City=%s Street=%s\n", e.City, e.Street)
	// 完整路径也可以
	fmt.Printf("完整: %s\n", e.Address.City)

	// ---- 4. tag 与 JSON ----
	u := User{ID: 1, Name: "Alice", Password: "secret"}
	b, _ := json.Marshal(u)
	fmt.Printf("JSON: %s\n", b) // Password 不出现,Email 也不出现

	// ---- 5. 空结构体 struct{}(0 字节) ----
	// 常用于 set 或作为 channel 的信号
	set := map[string]struct{}{}
	set["a"] = struct{}{}
	set["b"] = struct{}{}
	fmt.Printf("set size=%d\n", len(set))

	done := make(chan struct{})
	go func() {
		fmt.Println("goroutine 完成")
		close(done)
	}()
	<-done

	// ---- 6. 结构体比较 ----
	// 只要所有字段都可比较,结构体就可比较
	fmt.Println("p1 == p2:", p1 == p2)
	fmt.Println("Point{1,2} == Point{1,2}:", Point{1, 2} == Point{1, 2})
}
