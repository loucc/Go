// Package main 演示字符串处理。
//
// 学习要点:
//   - string 是不可变的 UTF-8 字节序列
//   - len(s) 返回字节数,不是字符数
//   - 遍历:按字节 for i := 0; i < len(s); i++;按 rune for _, r := range s
//   - []byte(s) 与 string(b) 都会复制底层数据
//   - 拼接大量字符串时用 strings.Builder,避免 O(n²)
//
// 运行:go run .
package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "Hello, 世界!"

	// ---- 1. 长度:字节 vs rune ----
	fmt.Printf("len(s)=%d(字节数)\n", len(s))
	fmt.Printf("utf8.RuneCountInString(s)=%d(字符数)\n", utf8.RuneCountInString(s))

	// ---- 2. 遍历(rune) ----
	for i, r := range s {
		fmt.Printf("  byte-index=%d rune=%c(U+%04X)\n", i, r, r)
	}

	// ---- 3. 截取(按字节!中文可能被截断) ----
	fmt.Println("s[:5] =", s[:5]) // "Hello"
	// 若要按字符截取,先转成 []rune
	runes := []rune(s)
	fmt.Println("runes[:7] =", string(runes[:7]))

	// ---- 4. strings 常用函数 ----
	fmt.Println(strings.ToUpper("hello"))
	fmt.Println(strings.Contains("golang", "go"))
	fmt.Println(strings.Split("a,b,c,d", ","))
	fmt.Println(strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Println(strings.Replace("foofoo", "foo", "bar", 1))
	fmt.Println(strings.ReplaceAll("foofoo", "foo", "bar"))
	fmt.Println(strings.TrimSpace("   hello   "))
	fmt.Println(strings.HasPrefix("golang.org", "go"))

	// ---- 5. strings.Builder(高效拼接) ----
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString("x")
	}
	fmt.Println("Builder:", sb.String())

	// ---- 6. 字符串与数字互转 ----
	n, err := strconv.Atoi("42")
	if err == nil {
		fmt.Printf("Atoi: %d\n", n)
	}
	fmt.Println("Itoa:", strconv.Itoa(100))

	f, _ := strconv.ParseFloat("3.14", 64)
	fmt.Printf("ParseFloat: %.2f\n", f)

	fmt.Println("FormatFloat:", strconv.FormatFloat(3.14159, 'f', 2, 64))

	// ---- 7. 转义与 raw string ----
	e := "line1\nline2"        // 双引号:支持转义
	r := `line1\nline2 \t raw` // 反引号:原样输出
	fmt.Println("---转义---")
	fmt.Println(e)
	fmt.Println("---raw---")
	fmt.Println(r)
}
