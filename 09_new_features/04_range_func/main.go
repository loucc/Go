// Package main 演示 Go 1.23:range over function(迭代器)。
//
// 学习要点:
//   - 可以对函数进行 range 循环
//   - 三种迭代器签名:
//     iter.Seq[V]:   func(yield func(V) bool)
//     iter.Seq2[K,V]: func(yield func(K, V) bool)
//   - yield 返回 false 时应停止迭代
//
// 运行:go run .
package main

import (
	"fmt"
	"iter"
)

// 一个简单的整数序列生成器
func Range(start, end int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return // 消费者提前退出
			}
		}
	}
}

// 一个 Key-Value 生成器
func Enumerate[T any](s []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}

func main() {
	// range over function
	for v := range Range(1, 6) {
		fmt.Println(v)
	}

	fmt.Println("---")

	// break 会传递给 yield
	for v := range Range(1, 100) {
		if v > 3 {
			break
		}
		fmt.Println(v)
	}

	fmt.Println("---")

	// K-V 迭代器
	words := []string{"go", "is", "great"}
	for i, w := range Enumerate(words) {
		fmt.Printf("%d: %s\n", i, w)
	}
}
