// Package main 演示切片(slice)——Go 中最重要的复合类型。
//
// 学习要点:
//   - slice 本质是 { ptr, len, cap } 的三元组
//   - append 可能触发扩容(重新分配底层数组)
//   - 多个 slice 可共享同一底层数组 → 修改可能相互影响
//   - 三索引 s[low:high:max] 限制新切片的 cap
//   - Go 1.21+ 的 slices 标准库提供了大量工具函数
//
// 运行:go run .
package main

import (
	"fmt"
	"slices"
)

func main() {
	// ---- 1. 创建 ----
	var s1 []int             // nil slice(len=0 cap=0),可以直接 append
	s2 := []int{1, 2, 3}     // 字面量
	s3 := make([]int, 3)     // [0 0 0]
	s4 := make([]int, 3, 10) // len=3 cap=10

	fmt.Printf("s1=%v (nil=%v) s2=%v s3=%v s4 len=%d cap=%d\n",
		s1, s1 == nil, s2, s3, len(s4), cap(s4))

	// ---- 2. append 的扩容 ----
	s := []int{}
	for i := 0; i < 8; i++ {
		s = append(s, i)
		fmt.Printf("append %d: len=%d cap=%d\n", i, len(s), cap(s))
	}

	// ---- 3. 共享底层数组的陷阱 ----
	base := []int{1, 2, 3, 4, 5}
	sub := base[1:3] // [2, 3]
	sub[0] = 99      // 会修改 base!
	fmt.Printf("修改 sub 后 base=%v sub=%v\n", base, sub)

	// ---- 4. 三索引限制 cap ----
	// s[low : high : max]  len = high-low, cap = max-low
	safe := base[1:3:3] // len=2 cap=2,再 append 会新分配,不影响 base
	safe = append(safe, 777)
	fmt.Printf("三索引 safe=%v base=%v (base 不受影响)\n", safe, base)

	// ---- 5. 删除元素(idx 位置) ----
	arr := []int{10, 20, 30, 40, 50}
	idx := 2
	arr = append(arr[:idx], arr[idx+1:]...) // 经典模式
	fmt.Printf("删除 idx=2 后: %v\n", arr)

	// 或用 slices.Delete(1.21+)
	arr2 := []int{10, 20, 30, 40, 50}
	arr2 = slices.Delete(arr2, 2, 3)
	fmt.Printf("slices.Delete: %v\n", arr2)

	// ---- 6. slices 标准库常用函数 ----
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	slices.Sort(nums)
	fmt.Printf("Sort: %v\n", nums)
	fmt.Printf("Contains 4: %v\n", slices.Contains(nums, 4))
	fmt.Printf("Index of 5: %v\n", slices.Index(nums, 5))
	fmt.Printf("Max: %d Min: %d\n", slices.Max(nums), slices.Min(nums))

	rev := []int{1, 2, 3, 4, 5}
	slices.Reverse(rev)
	fmt.Printf("Reverse: %v\n", rev)

	// Clone / Concat
	cloned := slices.Clone(rev)
	joined := slices.Concat(cloned, []int{100, 200})
	fmt.Printf("Concat: %v\n", joined)
}
