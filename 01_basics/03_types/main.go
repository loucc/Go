// Package main 演示 Go 的基本类型、零值与类型转换。
//
// 学习要点:
//   - 整数(有符号/无符号)、浮点、复数、bool、string
//   - 每种类型都有确定的零值
//   - Go 没有隐式类型转换,必须显式 T(x)
//   - byte = uint8;rune = int32(表示 Unicode 码点)
//
// 运行:go run .
package main

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func main() {
	// ---- 1. 整数类型 ----
	var i8 int8 = 127             // -128 ~ 127
	var u8 uint8 = 255            // 0 ~ 255
	var i64 int64 = math.MaxInt64 // 平台相关的 int 一般是 64 位
	fmt.Printf("int8=%d uint8=%d int64=%d\n", i8, u8, i64)

	// ---- 2. 浮点 ----
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793
	fmt.Printf("float32=%.6f float64=%.15f\n", f32, f64)

	// ---- 3. 布尔 ----
	var ok bool // 零值 false
	fmt.Printf("ok=%v (bool zero value)\n", ok)

	// ---- 4. 字符串与 rune ----
	s := "Hello, 世界"
	fmt.Printf("字节长度 len(s)=%d\n", len(s))                                       // 13(UTF-8)
	fmt.Printf("rune 长度 utf8.RuneCountInString=%d\n", utf8.RuneCountInString(s)) // 9

	// 用 for range 遍历字符串会按 rune 迭代
	for i, r := range s {
		fmt.Printf("  index=%d rune=%c(%U)\n", i, r, r)
	}

	// ---- 5. 类型转换(必须显式) ----
	x := 10   // int
	y := 3.14 // float64
	z := float64(x) + y
	fmt.Printf("z = %v (类型: %T)\n", z, z)

	// 字符串 <-> 字节 / rune
	b := []byte(s) // 复制底层
	r := []rune(s) // 复制并解码为 rune 序列
	fmt.Printf("[]byte len=%d, []rune len=%d\n", len(b), len(r))

	// ---- 6. 零值一览 ----
	var (
		zi  int
		zf  float64
		zs  string
		zb  bool
		zp  *int
		zsl []int
		zm  map[string]int
	)
	fmt.Printf("零值: int=%d float=%g string=%q bool=%v ptr=%v slice=%v map=%v\n",
		zi, zf, zs, zb, zp, zsl, zm)
}
