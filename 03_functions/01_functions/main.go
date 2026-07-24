// Package main 演示函数的高级用法。
//
// 学习要点:
//   - 多返回值(通常最后一个是 error)
//   - 命名返回值 + 裸 return(适度使用,提高可读性)
//   - 可变参数 ...T
//   - 函数是一等公民(可作参数、返回值、字段)
//   - 匿名函数与闭包
//
// 运行:go run .
package main

import (
	"errors"
	"fmt"
	"strings"
)

// 1. 多返回值
func divmod(a, b int) (int, int) {
	return a / b, a % b
}

// 2. 命名返回值 + 裸 return
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // 相当于 return x, y
}

// 3. 可变参数
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 4. 函数作为返回值(高阶函数)
func makeAdder(base int) func(int) int {
	return func(x int) int {
		return base + x
	}
}

// 5. 函数作为参数
func apply(nums []int, f func(int) int) []int {
	out := make([]int, len(nums))
	for i, n := range nums {
		out[i] = f(n)
	}
	return out
}

// 6. 返回 error
func safeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	// 多返回值
	q, r := divmod(17, 5)
	fmt.Printf("17 / 5 = %d, remainder %d\n", q, r)

	// 命名返回值
	x, y := split(17)
	fmt.Printf("split(17) = %d, %d\n", x, y)

	// 可变参数
	fmt.Println("sum(1,2,3,4,5) =", sum(1, 2, 3, 4, 5))

	// 传入切片时用 ...
	xs := []int{10, 20, 30}
	fmt.Println("sum(xs...) =", sum(xs...))

	// 闭包(每次调用 makeAdder 都产生独立的 base 变量)
	add10 := makeAdder(10)
	add100 := makeAdder(100)
	fmt.Println("add10(5) =", add10(5))
	fmt.Println("add100(5) =", add100(5))

	// 高阶函数(类似 map)
	squared := apply([]int{1, 2, 3, 4}, func(n int) int { return n * n })
	fmt.Println("squared =", squared)

	// 匿名函数直接调用
	greet := func(name string) string {
		return "Hello, " + strings.ToUpper(name)
	}
	fmt.Println(greet("go"))

	// error 处理
	if r, err := safeDivide(10, 0); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("result:", r)
	}
}
