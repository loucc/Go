// Package main 演示 map。
//
// 学习要点:
//   - map 是引用类型,zero value 是 nil(向 nil map 写入会 panic)
//   - 遍历顺序随机(每次运行都可能不同,防止依赖)
//   - v, ok := m[k] 用来判断键是否存在
//   - 非并发安全:并发读写必须加锁或用 sync.Map
//   - maps 包(Go 1.21+):Keys、Values、Clone、Copy、Equal
//
// 运行:go run .
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	// ---- 1. 创建 ----
	var m1 map[string]int  // nil map,不能写入!
	m2 := map[string]int{} // 空 map
	m3 := make(map[string]int, 10)
	m4 := map[string]int{"apple": 1, "banana": 2}

	fmt.Printf("m1 nil? %v m2=%v m3=%v m4=%v\n", m1 == nil, m2, m3, m4)

	// ---- 2. CRUD ----
	m4["cherry"] = 3
	m4["apple"] = 100 // 更新
	delete(m4, "banana")
	fmt.Printf("after CRUD: %v\n", m4)

	// ---- 3. 判断键是否存在 ----
	if v, ok := m4["apple"]; ok {
		fmt.Printf("apple exists, v=%d\n", v)
	}
	if _, ok := m4["banana"]; !ok {
		fmt.Println("banana not found")
	}

	// ---- 4. 遍历(顺序随机) ----
	for k, v := range m4 {
		fmt.Printf("  %s -> %d\n", k, v)
	}

	// ---- 5. 稳定遍历:先取 keys 排序 ----
	keys := slices.Sorted(maps.Keys(m4)) // 1.23+
	for _, k := range keys {
		fmt.Printf("[sorted] %s -> %d\n", k, m4[k])
	}

	// ---- 6. maps 标准库 ----
	c := maps.Clone(m4)
	fmt.Printf("Clone: %v\n", c)

	m5 := map[string]int{"cherry": 999, "date": 4}
	maps.Copy(c, m5) // 把 m5 合并到 c
	fmt.Printf("Copy 合并后: %v\n", c)

	fmt.Printf("Equal: %v\n", maps.Equal(m4, m4))

	// ---- 7. 复合值:map[string][]int ----
	groups := map[string][]int{}
	for _, n := range []int{1, 2, 3, 4, 5, 6} {
		key := "even"
		if n%2 == 1 {
			key = "odd"
		}
		groups[key] = append(groups[key], n)
	}
	fmt.Printf("groups: %v\n", groups)

	// ---- 8. 用作 set(map[T]struct{},内存最省) ----
	set := map[string]struct{}{}
	set["a"] = struct{}{}
	set["b"] = struct{}{}
	set["a"] = struct{}{} // 重复添加
	fmt.Printf("set size=%d\n", len(set))
}
