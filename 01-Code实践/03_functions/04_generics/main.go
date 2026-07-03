// Package main 演示 Go 1.18+ 的泛型。
//
// 学习要点:
//   - 泛型是"类型参数化",不是运行时机制
//   - 类型参数写在函数名/类型名后面的方括号里
//   - 约束(constraint)是一种特殊的接口,只用于类型参数
//   - any 是 interface{} 的别名
//   - comparable 表示支持 == 和 !=
//   - ~T 表示底层类型为 T 的所有类型(常用于自定义类型)
//   - Go 1.26:类型参数的约束里允许引用自己(F-bounded,见文末)
//   - 何时用:通用容器/算法。不要为一切类型都写泛型版本
//
// 运行:go run .
package main

import (
	"cmp"
	"fmt"
)

// 1. 最简单的泛型函数
func First[T any](s []T) T {
	var zero T
	if len(s) == 0 {
		return zero
	}
	return s[0]
}

// 2. 多个类型参数
func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

// 3. comparable 约束
func Contains[T comparable](s []T, x T) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

// 4. 自定义约束(接口)
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

// 5. cmp.Ordered:提供了内置的有序类型约束(Go 1.21+)
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 6. 泛型类型:栈
type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	n := len(s.data) - 1
	v := s.data[n]
	s.data = s.data[:n]
	return v, true
}

// 7. ~T 的意义:允许自定义底层类型也满足约束
type Celsius float64
type Fahrenheit float64

// 8. 自引用泛型(F-bounded)—— Go 1.26
//
// 约束里可以直接引用"自己",精确表达"我操作跟我同类型的值,并返回同类型"。
// 典型场景:数学/代数、Builder 链式 API、Clone、Merge。
// 用不上就别用,普通容器/算法不需要。

// Addable[A]:类型 A 必须提供 Add(A) A。
type Addable[A any] interface {
	Add(A) A
}

// AddAll 的返回类型和入参具体类型完全一致(编译期精确)。
func AddAll[A Addable[A]](xs ...A) A {
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

// Vec2 满足 Addable[Vec2]。
type Vec2 struct{ X, Y float64 }

func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Cloner[T]:Clone 必须返回同类型 T。
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
	fmt.Println("First([]int):", First([]int{10, 20, 30}))
	fmt.Println("First([]string):", First([]string{"a", "b"}))

	// Map 需要显式指定类型参数,或让编译器推断
	upper := Map([]int{1, 2, 3}, func(x int) string {
		return fmt.Sprintf("#%d", x*x)
	})
	fmt.Println("Map:", upper)

	fmt.Println("Contains:", Contains([]int{1, 2, 3}, 2))
	fmt.Println("Sum ints:", Sum([]int{1, 2, 3, 4}))
	fmt.Println("Sum floats:", Sum([]float64{1.1, 2.2, 3.3}))

	// 自定义类型也满足 Number(因为 ~float64)
	temps := []Celsius{20.5, 25.0, 18.3}
	fmt.Println("Sum Celsius:", Sum(temps))

	fmt.Println("Max:", Max(3, 5))
	fmt.Println("Max strings:", Max("apple", "banana"))

	// 泛型栈
	s := &Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	for {
		v, ok := s.Pop()
		if !ok {
			break
		}
		fmt.Println("pop:", v)
	}

	// 自引用泛型(1.26)
	v := AddAll(Vec2{1, 2}, Vec2{3, 4}, Vec2{10, 20})
	fmt.Printf("AddAll(Vec2...) = %+v\n", v)

	src := IntList{1, 2, 3}
	dst := CloneAny(src) // 返回类型精确保持 IntList,不需断言
	dst[0] = 99
	fmt.Println("CloneAny: src =", src, "dst =", dst)
}
