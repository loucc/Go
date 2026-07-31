// Package main 演示常见并发模式。
//
// 学习要点:
//   - Worker Pool:固定 goroutine 数处理任务队列(限流)
//   - Pipeline:多阶段流水线,每阶段用 channel 连接
//   - Fan-out / Fan-in:分发到多个 worker,再合并结果
//   - ErrGroup:并发任务组,任一失败则返回错误(生产用 x/sync/errgroup)
//   - Singleflight:相同 key 的并发请求合并为一次执行(防缓存击穿)
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

	fmt.Println("\n=== ErrGroup(手写版,生产用 golang.org/x/sync/errgroup) ===")
	errGroupDemo()

	fmt.Println("\n=== Singleflight(请求合并) ===")
	singleflightDemo()
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

// ---------- ErrGroup(手写简化版) ----------
//
// 生产环境请使用 golang.org/x/sync/errgroup,它提供:
//   - SetLimit(n) 限制并发数
//   - 上下文取消:任一 goroutine 返回 error 时自动取消其余
//   - Go 1.25+ 可用 sync.WaitGroup.Go 简化写法
//
// 这里用标准库手写一个简化版,帮助理解原理。
type errGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

func (g *errGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if e := f(); e != nil {
			g.errOnce.Do(func() { g.err = e })
		}
	}()
}

func (g *errGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func errGroupDemo() {
	var g errGroup

	for i := 1; i <= 5; i++ {
		id := i
		g.Go(func() error {
			time.Sleep(time.Duration(id*20) * time.Millisecond)
			fmt.Printf("  task %d done\n", id)
			return nil
		})
	}

	// 模拟一个失败的任务
	g.Go(func() error {
		time.Sleep(30 * time.Millisecond)
		fmt.Println("  task-fail: simulating error")
		return fmt.Errorf("task-fail: something went wrong")
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("  errgroup 返回首个错误: %v\n", err)
	}
}

// ---------- Singleflight(请求合并) ----------
//
// 场景:缓存击穿 —— 同一个 key 同时有大量并发请求,
// 如果缓存 miss,它们会同时打到下游(如数据库)。
// Singleflight 保证同一 key 的并发请求中只有一个真正执行,
// 其余等待结果共享。
//
// 生产环境请使用 golang.org/x/sync/singleflight。
type singleflight struct {
	mu sync.Mutex
	m  map[string]*call
}

type call struct {
	wg  sync.WaitGroup
	val string
	err error
}

func (sf *singleflight) Do(key string, fn func() (string, error)) (string, error) {
	sf.mu.Lock()
	if sf.m == nil {
		sf.m = make(map[string]*call)
	}
	if c, ok := sf.m[key]; ok {
		sf.mu.Unlock()
		c.wg.Wait() // 已有相同 key 的请求在执行,等待它
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	sf.m[key] = c
	sf.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	sf.mu.Lock()
	delete(sf.m, key) // 完成后删除,下次请求重新执行
	sf.mu.Unlock()

	return c.val, c.err
}

func singleflightDemo() {
	var sf singleflight

	// 模拟 10 个并发请求同一个 key,但只有 1 个真正执行 fetch
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val, err := sf.Do("user:1", func() (string, error) {
				fmt.Println("  [实际执行] fetch user:1 from database (只会出现一次)")
				time.Sleep(50 * time.Millisecond)
				return "Alice", nil
			})
			if err != nil {
				fmt.Printf("  worker %d: error %v\n", id, err)
				return
			}
			fmt.Printf("  worker %d: got %q\n", id, val)
		}(i)
	}
	wg.Wait()
}
