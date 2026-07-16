// Package main 演示 select 多路复用。
//
// 学习要点:
//   - select 会等待多个 channel 中第一个就绪的分支
//   - 所有分支都阻塞时,select 也阻塞;有 default 则不阻塞
//   - 多个分支同时就绪时,随机选择一个(防止饥饿)
//   - nil channel 永远阻塞,可用来动态"关闭"某个分支
//   - 与 time.After / context.Done() 结合做超时/取消
//
// 运行:go run .
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ---- 1. 超时模式 ----
	ch := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "结果"
	}()

	select {
	case v := <-ch:
		fmt.Println("收到:", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时!")
	}

	// ---- 2. 非阻塞发送/接收(default) ----
	quick := make(chan int, 1)
	select {
	case quick <- 100:
		fmt.Println("成功发送 100")
	default:
		fmt.Println("channel 满,跳过")
	}
	select {
	case v := <-quick:
		fmt.Println("成功接收:", v)
	default:
		fmt.Println("没数据")
	}

	// ---- 3. context 取消 ----
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	slow := make(chan int)
	go func() {
		time.Sleep(500 * time.Millisecond)
		slow <- 999
	}()

	select {
	case v := <-slow:
		fmt.Println("拿到 slow:", v)
	case <-ctx.Done():
		fmt.Println("context 取消:", ctx.Err())
	}

	// ---- 4. 多路合并(fan-in) ----
	a := gen("A", 3)
	b := gen("B", 3)
	for range 6 {
		select {
		case v := <-a:
			fmt.Println("来自 A:", v)
		case v := <-b:
			fmt.Println("来自 B:", v)
		}
	}
}

func gen(tag string, n int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := range n {
			ch <- fmt.Sprintf("%s-%d", tag, i)
			time.Sleep(30 * time.Millisecond)
		}
	}()
	return ch
}
