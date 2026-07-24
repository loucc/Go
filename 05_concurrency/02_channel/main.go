// Package main 演示 channel。
//
// 学习要点:
//   - channel 是并发原语,用来 goroutine 间通信
//   - 无缓冲 channel:同步的,发送与接收必须同时准备好
//   - 有缓冲 channel:异步的,缓冲未满/未空前不阻塞
//   - close(ch) 后仍可读完缓冲,但读到零值(第二个返回值为 false)
//   - 向已 close 的 channel 写会 panic
//   - 单向 channel 用于类型收窄:chan<- T 只写,<-chan T 只读
//
// 运行:go run .
package main

import (
	"fmt"
	"time"
)

func main() {
	// ---- 1. 无缓冲 channel:同步 ----
	ch := make(chan int)
	go func() {
		fmt.Println("发送前")
		ch <- 42 // 阻塞直到有接收方
		fmt.Println("发送后")
	}()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("main 接收:", <-ch)

	// ---- 2. 有缓冲 channel:异步 ----
	buf := make(chan int, 3)
	buf <- 1
	buf <- 2
	buf <- 3
	// buf <- 4 // 会阻塞!缓冲已满
	fmt.Printf("buf len=%d cap=%d\n", len(buf), cap(buf))
	close(buf)

	// close 后依然可以读完
	for v := range buf {
		fmt.Println("读取:", v)
	}

	// ---- 3. 单向 channel:限制方向,增强类型安全 ----
	nums := make(chan int, 5)
	go producer(nums)
	consumer(nums)

	// ---- 4. 双返回值判断关闭 ----
	ch2 := make(chan int)
	close(ch2)
	if v, ok := <-ch2; !ok {
		fmt.Printf("channel 已关闭,v=%d(零值),ok=%v\n", v, ok)
	}
}

// 只写 channel:意图更清晰
func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i * i
	}
	close(ch)
}

// 只读 channel
func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println("consumer 收到:", v)
	}
}
