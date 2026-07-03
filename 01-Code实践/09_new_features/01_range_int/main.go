// Package main 演示 Go 1.22:for range 整数迭代。
//
// 之前:for i := 0; i < n; i++
// 现在:for i := range n
package main

import "fmt"

func main() {
	// 直接对整数迭代:0..4
	for i := range 5 {
		fmt.Println(i)
	}

	// 只关心次数,不需要变量
	count := 0
	for range 10 {
		count++
	}
	fmt.Println("count =", count)
}
