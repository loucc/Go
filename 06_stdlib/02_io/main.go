// Package main 演示 io 包与 os 文件操作。
//
// 学习要点:
//   - io.Reader / io.Writer 是最重要的两个接口
//   - bufio 提供带缓冲的读写(减少系统调用)
//   - os.ReadFile / os.WriteFile 是一次性读写整个文件的便捷方法
//   - 读写完记得 defer Close
//
// 运行:go run .
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// ---- 1. 一次性读写 ----
	tmp := os.TempDir() + "/go_learning_demo.txt"
	content := "第一行\n第二行\n第三行\n"
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		panic(err)
	}
	defer os.Remove(tmp)

	data, err := os.ReadFile(tmp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("读到 %d 字节:\n%s", len(data), data)

	// ---- 2. 按行读取(bufio.Scanner) ----
	f, err := os.Open(tmp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("第 %d 行: %s\n", i, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scan err:", err)
	}

	// ---- 3. Reader / Writer 抽象 ----
	// 任何实现 Read 的类型都是 io.Reader,strings.Reader 就是一个
	r := strings.NewReader("Hello, io!")
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Printf("读到 %d 字节: %q\n", n, string(buf[:n]))
		}
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
	}

	// ---- 4. io.Copy(高效流式拷贝) ----
	src := strings.NewReader("from src\n")
	// os.Stdout 实现了 io.Writer
	n, _ := io.Copy(os.Stdout, src)
	fmt.Printf("(拷贝 %d 字节)\n", n)
}
