// Package main 演示 Go 1.26:自引用泛型 (F-bounded polymorphism)。
//
// 1.26 之前:类型参数不能在自己的约束里引用"自身"。
//
//	type Adder[A Adder[A]] interface { Add(A) A }  // 编译错误
//	只能用普通接口 + 类型断言,或者多余的辅助类型参数绕开。
//
// 1.26 之后:允许。可以精确表达"我操作跟我同类型的值,并返回同类型"。
//
// 用途:
//   - 数学/代数类型:Add、Compare、Merge,返回同种类型
//   - Builder / 链式 API:每一步返回具体 builder 类型
//   - Clone:Clone() 必须返回同种具体类型
//
// 别滥用:普通容器/算法很少需要。只有真的要"我 = 我处理的东西"时才用。
package main

import "fmt"

// Adder 约束:类型 A 必须提供一个 Add(A) A 方法,
// 也就是"操作 A 类型并返回 A 类型"。
type Adder[A any] interface {
	Add(A) A
}

// Sum 只对满足自引用约束的类型可用,返回类型永远精确。
func Sum[A Adder[A]](xs ...A) A {
	var acc A
	for i, v := range xs {
		if i == 0 {
			acc = v
			continue
		}
		acc = acc.Add(v)
	}
	return acc
}

// Vec2 是一个满足 Adder[Vec2] 的具体类型。
type Vec2 struct{ X, Y float64 }

func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Money 也满足 Adder[Money]。
type Money int64

func (a Money) Add(b Money) Money { return a + b }

// Cloner:典型 F-bounded 用法,Clone 必须返回同类型。
type Cloner[T any] interface {
	Clone() T
}

func CloneAny[T Cloner[T]](v T) T { return v.Clone() }

type IntList []int

func (l IntList) Clone() IntList {
	c := make(IntList, len(l))
	copy(c, l)
	return c
}

func main() {
	// 数学:结果类型自动是 Vec2,不需要断言
	v := Sum(Vec2{1, 2}, Vec2{3, 4}, Vec2{10, 20})
	fmt.Printf("Sum(Vec2...) = %+v\n", v)

	m := Sum[Money](1, 2, 3, 4)
	fmt.Println("Sum(Money...) =", m)

	// Clone:返回类型精确保持为 IntList
	src := IntList{1, 2, 3}
	dst := CloneAny(src)
	dst[0] = 99
	fmt.Println("src =", src, "dst =", dst)
}
