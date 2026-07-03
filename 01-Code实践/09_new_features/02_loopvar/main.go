// Package main 演示 Go 1.22 循环变量作用域的变化。
//
// 在 1.22 之前:
//
//	for i := 0; i < 3; i++ { go func() { fmt.Println(i) }() }
//	会输出三个 3,因为所有 goroutine 共享同一个 i。
//
// 从 1.22 开始:
//
//	循环每次迭代都产生新的 i,goroutine 捕获的是不同的变量。
//	输出是 0/1/2(顺序不定)。
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() { // 1.22+ 中每次迭代 i 都是新变量
			defer wg.Done()
			fmt.Println("i =", i)
		}()
	}
	wg.Wait()

	// range 循环也一样
	nums := []int{10, 20, 30}
	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("n =", n)
		}()
	}
	wg.Wait()
}
