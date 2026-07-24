// Package main 演示 unsafe 包。
//
// ⚠️  unsafe 顾名思义"不安全",绕过 Go 的类型系统。
//
//	仅在:高性能场景、与 C 互操作、了解内存布局时使用。
//
// 学习要点:
//   - unsafe.Pointer:任意指针间的桥梁
//   - unsafe.Sizeof / Alignof / Offsetof:内存布局工具
//   - unsafe.Slice / String(1.17+ / 1.20+):零拷贝构造 slice / string
//
// 运行:go run .
package main

import (
	"fmt"
	"unsafe"
)

type Point struct {
	X int32 // offset 0
	Y int32 // offset 4
	Z int64 // offset 8(考虑到 8 字节对齐)
}

func main() {
	// ---- 1. 内存布局 ----
	p := Point{}
	fmt.Printf("Sizeof(Point) = %d\n", unsafe.Sizeof(p))
	fmt.Printf("Alignof(p.Y)  = %d\n", unsafe.Alignof(p.Y))
	fmt.Printf("Offsetof(Y)   = %d\n", unsafe.Offsetof(p.Y))
	fmt.Printf("Offsetof(Z)   = %d\n", unsafe.Offsetof(p.Z))

	// ---- 2. []byte 与 string 零拷贝转换(实验性!) ----
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	// 从 []byte 得到 string(不复制)
	s := unsafe.String(unsafe.SliceData(b), len(b))
	fmt.Println("string:", s)

	// ⚠️ 之后修改 b 会同时"修改"这个 string,违反 string 不可变约定!
	// 只在你 100% 确定不修改时才这样做。

	// ---- 3. 通过 uintptr 做指针算术(危险,不推荐) ----
	// Go 的 GC 可能移动对象,uintptr 只是"数字",不追踪
	// 因此下面的写法在实际业务代码里应避免
	x := int64(0x1234)
	pp := unsafe.Pointer(&x)
	pi := (*int32)(pp) // 强制把 *int64 当 *int32 看
	fmt.Printf("低 32 位 = 0x%x\n", *pi)

	// ---- 4. 结构体字段"顺序访问"(用 Offsetof) ----
	base := unsafe.Pointer(&p)
	xp := (*int32)(unsafe.Pointer(uintptr(base) + unsafe.Offsetof(p.X)))
	yp := (*int32)(unsafe.Pointer(uintptr(base) + unsafe.Offsetof(p.Y)))
	*xp = 100
	*yp = 200
	fmt.Printf("修改后 p = %+v\n", p)
}
