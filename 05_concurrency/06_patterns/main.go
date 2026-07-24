// Package main 演示常见并发模式:worker pool、pipeline、fan-in/fan-out。
//
// 学习要点:
//   - Worker Pool:固定 goroutine 数处理任务队列(限流)
//   - Pipeline:多阶段流水线,每阶段用 channel 连接
//   - Fan-out / Fan-in:分发到多个 worker,再合并结果
//   - 关闭 channel 的原则:由发送方关闭(明确唯一)
//
// 运行:go run .
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Worker Pool ===")
	workerPool()

	fmt.Println("\n=== Pipeline ===")
	pipeline()

	fmt.Println("\n=== Fan-out / Fan-in ===")
	fanOutFanIn()
}

// ---------- Worker Pool ----------
func workerPool() {
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// 启动 3 个 worker
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				time.Sleep(20 * time.Millisecond)
				results <- j * j
				fmt.Printf("  worker %d 处理 job=%d\n", id, j)
			}
		}(w)
	}

	// 发任务
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // 关闭任务队列,worker 从 range 循环退出

	// 等 worker 全部退出后关闭 results
	go func() {
		wg.Wait()
		close(results)
	}()

	// 消费结果
	for r := range results {
		fmt.Printf("  结果: %d\n", r)
	}
}

// ---------- Pipeline ----------
func pipeline() {
	// stage 1: 生成 1..5
	nums := make(chan int)
	go func() {
		defer close(nums)
		for i := 1; i <= 5; i++ {
			nums <- i
		}
	}()

	// stage 2: 平方
	squared := make(chan int)
	go func() {
		defer close(squared)
		for n := range nums {
			squared <- n * n
		}
	}()

	// stage 3: 打印
	for v := range squared {
		fmt.Println("  =>", v)
	}
}

// ---------- Fan-out / Fan-in ----------
func fanOutFanIn() {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			input <- i
		}
	}()

	// Fan-out:多个 worker 并行处理
	worker := func(id int, in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for v := range in {
				time.Sleep(10 * time.Millisecond)
				out <- fmt.Sprintf("worker-%d 处理 %d", id, v)
			}
		}()
		return out
	}
	c1 := worker(1, input)
	c2 := worker(2, input)
	c3 := worker(3, input)

	// Fan-in:合并多个 channel
	merged := fanIn(c1, c2, c3)
	for msg := range merged {
		fmt.Println("  ", msg)
	}
}

func fanIn(chs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
