// Package main 演示 Go 1.23 的 iter/slices/maps 迭代器 API。
//
// - slices.All(s):   返回 iter.Seq2[int, T]
// - slices.Values(s):返回 iter.Seq[T]
// - maps.Keys(m):    返回 iter.Seq[K]
// - maps.Values(m):  返回 iter.Seq[V]
// - slices.Collect / slices.Sorted 把迭代器收集回 slice
//
// 运行:go run .
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	s := []string{"go", "rust", "python"}

	// slices.All 返回 (index, value)
	for i, v := range slices.All(s) {
		fmt.Printf("%d -> %s\n", i, v)
	}

	// slices.Values 只返回 value
	for v := range slices.Values(s) {
		fmt.Println("v:", v)
	}

	// maps 迭代器
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := slices.Sorted(maps.Keys(m)) // 收集 + 排序
	fmt.Println("keys:", keys)

	// Collect 把任意 iter.Seq 收集回 slice
	nums := slices.Collect(func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			if !yield(i * i) {
				return
			}
		}
	})
	fmt.Println("collected:", nums)
}
