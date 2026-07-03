// Package main 演示数组。
//
// 学习要点:
//   - 数组是值类型,长度是类型的一部分:[3]int 与 [4]int 不同
//   - 数组赋值/传参会复制整个数组(性能陷阱!)
//   - 生产中很少直接用数组,更多用切片
//
// 运行:go run .
package main

import "fmt"

func main() {
	// 声明与初始化
	var a [3]int                // [0 0 0]
	b := [3]int{1, 2, 3}        // 完全指定
	c := [...]int{10, 20, 30}   // 编译器推断长度
	d := [5]int{1: 100, 3: 300} // 指定下标 [0 100 0 300 0]

	fmt.Printf("a=%v b=%v c=%v d=%v\n", a, b, c, d)

	// 二维数组
	var grid [2][3]int
	grid[0] = [3]int{1, 2, 3}
	grid[1] = [3]int{4, 5, 6}
	fmt.Printf("grid=%v\n", grid)

	// 数组比较(必须同类型)
	x := [3]int{1, 2, 3}
	y := [3]int{1, 2, 3}
	fmt.Printf("x == y: %v\n", x == y)

	// 值语义:传参会复制
	modifyArray(b)
	fmt.Printf("after modifyArray, b=%v (未改变,因为传值)\n", b)

	// 若想让函数修改原数组,传指针
	modifyArrayPtr(&b)
	fmt.Printf("after modifyArrayPtr, b=%v\n", b)
}

func modifyArray(a [3]int) {
	a[0] = 999
}

func modifyArrayPtr(a *[3]int) {
	a[0] = 999
}
