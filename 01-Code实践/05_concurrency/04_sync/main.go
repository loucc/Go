// Package main 演示 sync 包核心工具。
//
// 学习要点:
//   - sync.Mutex / RWMutex:互斥/读写锁
//   - sync.WaitGroup:等待一组 goroutine 完成
//   - sync.Once:仅执行一次(单例、懒加载)
//   - sync.Pool:临时对象复用池,减轻 GC 压力
//   - sync.Map:读多写少的并发 map
//   - sync/atomic:无锁原子操作,性能最好
//   - 用 go run -race 检测数据竞争
//
// 运行:go run .
//
//	go run -race .   # 开启数据竞争检测
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ---- 1. Mutex 保护共享状态 ----
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// ---- 2. RWMutex(读远多于写时更快) ----
type Cache struct {
	mu sync.RWMutex
	m  map[string]string
}

func (c *Cache) Get(k string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.m[k]
}

func (c *Cache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = v
}

// ---- 3. sync.Once(单例初始化) ----
var (
	once     sync.Once
	instance *Cache
)

func GetCache() *Cache {
	once.Do(func() {
		fmt.Println("[Once] 初始化 cache")
		instance = &Cache{m: map[string]string{}}
	})
	return instance
}

func main() {
	// Mutex 演示
	c := &SafeCounter{}
	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() {
			c.Inc()
		})
	}
	wg.Wait()
	fmt.Println("Mutex count:", c.count)

	// atomic:比 Mutex 更快
	var atomicCounter int64
	wg = sync.WaitGroup{}
	for range 1000 {
		wg.Go(func() {
			atomic.AddInt64(&atomicCounter, 1)
		})
	}
	wg.Wait()
	fmt.Println("Atomic count:", atomic.LoadInt64(&atomicCounter))

	// Once
	for i := range 3 {
		GetCache().Set(fmt.Sprintf("key%d", i), "v")
	}
	fmt.Println("cache len:", len(GetCache().m))

	// sync.Pool:临时对象池
	pool := sync.Pool{
		New: func() any {
			return make([]byte, 1024)
		},
	}
	buf := pool.Get().([]byte)
	// 用完放回
	pool.Put(buf)
	fmt.Println("pool 演示完成,buf len =", len(buf))

	// sync.Map:并发安全的 map
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)
	if v, ok := m.Load("a"); ok {
		fmt.Println("sync.Map a =", v)
	}
	m.Range(func(k, v any) bool {
		fmt.Printf("  %v -> %v\n", k, v)
		return true // 返回 false 停止遍历
	})
}
