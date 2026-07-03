// Package main 演示 Go 的内建函数。
//
// 学习要点:
//   - len / cap:长度与容量(数组/切片/字符串/channel/map)
//   - make:分配 slice / map / chan 的运行时结构
//   - new:分配零值并返回指针
//   - append / copy:切片操作
//   - delete:从 map 删除键
//   - close:关闭 channel
//   - min / max / clear:Go 1.21+ 新增
//   - panic / recover:异常机制(慎用)
//
// 运行:go run .
package main

import "fmt"

func main() {
	// ---- 1. len 与 cap ----
	s := make([]int, 3, 10) // len=3 cap=10
	fmt.Printf("len(s)=%d cap(s)=%d\n", len(s), cap(s))

	m := map[string]int{"a": 1, "b": 2}
	fmt.Printf("len(m)=%d\n", len(m))

	// ---- 2. make vs new ----
	// make 只用于 slice/map/channel,返回类型 T
	sl := make([]int, 5)
	// new 返回 *T,分配零值内存
	pi := new(int) // pi 指向 int 零值 0
	*pi = 42
	fmt.Printf("slice=%v, *pi=%d\n", sl, *pi)

	// ---- 3. append 与 copy ----
	a := []int{1, 2, 3}
	a = append(a, 4, 5)
	fmt.Printf("append 后: %v\n", a)

	b := make([]int, len(a))
	copy(b, a)
	fmt.Printf("copy 后: %v\n", b)

	// ---- 4. delete ----
	delete(m, "a")
	fmt.Printf("delete 后: %v\n", m)

	// ---- 5. close(channel) ----
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	close(ch)
	for v := range ch { // 关闭后依然可以读完缓冲区
		fmt.Printf("channel value: %d\n", v)
	}

	// ---- 6. min / max / clear(Go 1.21+) ----
	fmt.Println("min(3,7,1) =", min(3, 7, 1))
	fmt.Println("max(3,7,1) =", max(3, 7, 1))

	nums := []int{1, 2, 3, 4, 5}
	clear(nums) // 将 slice 元素全部置零值(注意:不改变 len)
	fmt.Printf("clear 后: %v\n", nums)

	dict := map[string]int{"x": 1, "y": 2}
	clear(dict) // 清空 map 的所有键
	fmt.Printf("clear map 后: %v (len=%d)\n", dict, len(dict))

	// ---- 7. panic / recover ----
	safeDivide(10, 0)
	fmt.Println("程序继续执行 ✔")
}

func safeDivide(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover from panic: %v\n", r)
		}
	}()
	if b == 0 {
		panic("division by zero")
	}
	fmt.Println(a / b)
}
