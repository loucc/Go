// Package main 演示控制流:if / for / switch / goto。
//
// 学习要点:
//   - Go 只有 for 一种循环关键字
//   - if / switch 支持初始化子句
//   - switch 默认无 fallthrough(与 C 相反),按需显式声明
//   - Go 1.22+:for range 整数、循环变量每次迭代新建
//
// 运行:go run .
package main

import "fmt"

func main() {
	// ---- 1. if / else / 带初始化 ----
	if n := 42; n%2 == 0 {
		fmt.Printf("%d 是偶数\n", n)
	} else {
		fmt.Printf("%d 是奇数\n", n)
	}

	// ---- 2. for 的四种写法 ----
	// (a) C 风格三段式
	sum := 0
	for i := 1; i <= 10; i++ {
		sum += i
	}
	fmt.Println("1..10 sum =", sum)

	// (b) while 式
	n := 8
	for n > 1 {
		n /= 2
	}
	fmt.Println("n =", n)

	// (c) 无限循环 + break
	count := 0
	for {
		count++
		if count == 3 {
			break
		}
	}
	fmt.Println("count =", count)

	// (d) for range(1.22+ 也支持对整数迭代)
	for i := range 5 { // 0,1,2,3,4
		fmt.Printf("range %d ", i)
	}
	fmt.Println()

	// ---- 3. switch(默认无 fallthrough)----
	day := 3
	switch day {
	case 1, 2, 3, 4, 5:
		fmt.Println("工作日")
	case 6, 7:
		fmt.Println("周末")
	default:
		fmt.Println("无效")
	}

	// 表达式 switch(无 tag,类似 if-else 链)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B")
	case score >= 60:
		fmt.Println("C")
	default:
		fmt.Println("D")
	}

	// 类型 switch(常见于处理 interface{})
	describe(42)
	describe("hello")
	describe(3.14)
	describe(true)

	// ---- 4. 标签 + break/continue(跳出多重循环) ----
outer:
	for i := range 3 {
		for j := range 3 {
			if i*j > 3 {
				fmt.Printf("break outer at i=%d j=%d\n", i, j)
				break outer
			}
		}
	}
}

func describe(i any) {
	switch v := i.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %q\n", v)
	case float64:
		fmt.Printf("float64: %g\n", v)
	default:
		fmt.Printf("其他类型 %T: %v\n", v, v)
	}
}
