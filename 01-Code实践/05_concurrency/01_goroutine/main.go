// Package main 演示 goroutine。
//
// 学习要点:
//   - go f() 启动一个 goroutine,几乎不阻塞
//   - 初始栈只有几 KB,按需扩容,可轻松开百万级
//   - 主 goroutine 退出,程序直接结束(不等其他 goroutine)
//   - 用 WaitGroup / channel 等待 goroutine 完成
//   - GOMAXPROCS 控制并行的 OS 线程数(Go 1.25+ 会读取容器 cgroup 限制)
//
// 运行:go run .
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Printf("GOMAXPROCS=%d NumGoroutine=%d\n",
		runtime.GOMAXPROCS(0), runtime.NumGoroutine())

	// ---- 1. 最简单的 goroutine ----
	go func() {
		fmt.Println("hello from goroutine")
	}()
	// 如果这里不 Sleep,主 goroutine 结束程序就退出了
	time.Sleep(50 * time.Millisecond)

	// ---- 2. 用 WaitGroup 等待完成 ----
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("worker %d done\n", id)
		}(i) // ← 传 i 进去,防止循环变量陷阱(1.22 前尤其重要)
	}
	wg.Wait()
	fmt.Println("all workers done")

	// ---- 3. 观察 goroutine 数量 ----
	fmt.Printf("最终 NumGoroutine=%d\n", runtime.NumGoroutine())
}
