// Package main 演示 Go 1.24:泛型类型别名。
//
// 1.24 之前:type A = B 只能是具体类型的别名,不能带类型参数
// 1.24 之后:type Set[T comparable] = map[T]struct{}
package main

import "fmt"

// 泛型类型别名(1.24+)
type Set[T comparable] = map[T]struct{}

func newSet[T comparable]() Set[T] { return Set[T]{} }

func add[T comparable](s Set[T], v T) { s[v] = struct{}{} }

func has[T comparable](s Set[T], v T) bool {
	_, ok := s[v]
	return ok
}

func main() {
	s := newSet[string]()
	add(s, "a")
	add(s, "b")
	add(s, "a")
	fmt.Println("size =", len(s))
	fmt.Println("has(a) =", has(s, "a"))
	fmt.Println("has(z) =", has(s, "z"))
}
