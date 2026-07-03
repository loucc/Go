// Package main 演示 defer。
//
// 学习要点:
//   - defer 在函数返回前执行(LIFO 后进先出)
//   - defer 的参数在 defer 语句执行时就求值,不是执行时
//   - 常用于:关闭资源、解锁、recover panic
//   - 循环中 defer 需谨慎(会累积到函数结束才执行)
//   - Go 1.14 起 defer 性能已经非常接近直接调用
//
// 运行:go run .
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== 1. LIFO 顺序 ===")
	demoOrder()

	fmt.Println("\n=== 2. 参数即时求值 ===")
	demoArgEval()

	fmt.Println("\n=== 3. defer 修改命名返回值 ===")
	fmt.Println("result =", demoModifyReturn())

	fmt.Println("\n=== 4. 资源清理模式 ===")
	if err := readFile("/etc/hostname"); err != nil {
		fmt.Println("err:", err)
	}
}

func demoOrder() {
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("main body")
	// 输出:main body, defer 3, defer 2, defer 1
}

func demoArgEval() {
	i := 10
	defer fmt.Println("defer i =", i) // 这里 i 立即求值为 10
	i = 20
	fmt.Println("current i =", i)
	// 输出:current i = 20 → defer i = 10
}

func demoModifyReturn() (result int) {
	defer func() {
		result *= 2 // 可以修改命名返回值!
	}()
	return 5 // 先设 result=5,再执行 defer,最后真正返回 10
}

func readFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close() // 无论后面怎么退出都会执行

	buf := make([]byte, 128)
	n, err := f.Read(buf)
	if err != nil {
		return err
	}
	fmt.Printf("读入 %d 字节: %s", n, string(buf[:n]))
	return nil
}
