// Package main 演示 cmp 包(Go 1.21+)与 slices/sort 的配合。
//
// cmp 包提供两个核心函数:
//   - cmp.Compare(a, b):返回 -1 / 0 / 1,满足 cmp.Ordered 约束的类型
//   - cmp.Less(a, b):返回 a < b
//
// cmp.Or 是另一个实用函数:
//   - cmp.Or(a, b, ...):返回第一个非零值(类似 SQL 的 COALESCE)
//
// 典型场景:
//   - 自定义排序:sort.Slice + cmp.Compare
//   - 多字段排序:先按字段 A 比,相等再按字段 B 比
//   - 默认值选择:cmp.Or(userInput, defaultValue)
//
// 运行:go run .
package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Student struct {
	Name  string
	Score int
	Age   int
}

func main() {
	// ---- 1. cmp.Compare 基础 ----
	fmt.Println("Compare(3, 5):", cmp.Compare(3, 5))       // -1
	fmt.Println("Compare(5, 3):", cmp.Compare(5, 3))       // 1
	fmt.Println("Compare(3, 3):", cmp.Compare(3, 3))       // 0
	fmt.Println("Compare(\"a\", \"b\"):", cmp.Compare("a", "b")) // -1
	fmt.Println("Less(3, 5):", cmp.Less(3, 5))             // true

	// ---- 2. 自定义排序:按分数降序,分数相同按年龄升序 ----
	students := []Student{
		{"Alice", 90, 20},
		{"Bob", 85, 22},
		{"Carol", 90, 19},
		{"Dave", 78, 21},
		{"Eve", 90, 20},
	}

	slices.SortFunc(students, func(a, b Student) int {
		// 先按分数降序(用 b 比 a)
		if c := cmp.Compare(b.Score, a.Score); c != 0 {
			return c
		}
		// 分数相同,按年龄升序
		return cmp.Compare(a.Age, b.Age)
	})

	fmt.Println("\n排序后(分数降序 → 年龄升序):")
	for _, s := range students {
		fmt.Printf("  %-6s score=%d age=%d\n", s.Name, s.Score, s.Age)
	}

	// ---- 3. slices.BinarySearch 用 cmp.Ordered ----
	nums := []int{1, 3, 5, 7, 9, 11, 13}
	pos, found := slices.BinarySearch(nums, 7)
	fmt.Printf("\nBinarySearch(7): pos=%d found=%v\n", pos, found)
	pos, found = slices.BinarySearch(nums, 8)
	fmt.Printf("BinarySearch(8): pos=%d found=%v (应插入位置)\n", pos, found)

	// ---- 4. cmp.Or:返回第一个非零值(Go 1.22+) ----
	// 类似 SQL 的 COALESCE / JavaScript 的 ??
	fmt.Println()
	fmt.Println("cmp.Or(\"\", \"default\"):", cmp.Or("", "default"))
	fmt.Println("cmp.Or(0, 42):", cmp.Or(0, 42))
	fmt.Println("cmp.Or(nil_err, err):", cmp.Or[error](nil, fmt.Errorf("fallback error")))

	// 实用场景:函数参数默认值
	port := 0
	fmt.Printf("port = cmp.Or(port, 8080) → %d\n", cmp.Or(port, 8080))

	// ---- 5. 用 cmp.Ordered 写泛型排序工具 ----
	fmt.Println()
	sorted := sortBy([]string{"banana", "apple", "cherry"})
	fmt.Println("sortBy strings:", sorted)

	sortedInts := sortBy([]int{5, 2, 8, 1, 9})
	fmt.Println("sortBy ints:", sortedInts)
}

// sortBy 用 cmp.Ordered 约束实现通用排序。
// 等价于 slices.Clone + slices.Sort,这里演示原理。
func sortBy[T cmp.Ordered](s []T) []T {
	out := slices.Clone(s)
	slices.SortFunc(out, cmp.Compare)
	return out
}
