// Package main 演示 Go 程序的最小结构。
//
// 学习要点:
//  1. 每个可执行程序必须有且只有一个 main 包
//  2. import 用于导入其他包,可分行或分组
//  3. main 函数是程序入口,没有参数也没有返回值
//  4. init 函数会在 main 之前执行,可以有多个
//
// 运行:go run .
package main

import "fmt"

// init 会在包被加载时自动执行,先于 main
// 同一个包内可以有多个 init(按文件字典序执行)。
func init() {
	fmt.Println("[init] 包初始化完成")
}

func main() {
	fmt.Println("Hello, Go 1.26.4!")
	fmt.Printf("这是我的第一个 Go 程序,当前包名: %s\n", "main")
}
