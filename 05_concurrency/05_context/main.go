// Package main 演示 context 包。
//
// 学习要点:
//   - context 用于:取消信号、超时、截止时间、请求范围的值
//   - 约定:作为函数的第一个参数,变量名 ctx
//   - 不要把 context 存入结构体,也不要传 nil(用 context.TODO)
//   - 不要用 context.WithValue 传业务参数,只传"横切"数据(如 traceID)
//   - 一旦父 context 取消,所有子 context 立即取消(级联)
//
// 运行:go run .
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	// ---- 1. WithCancel:手动取消 ----
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx, "手动取消 worker")
	time.Sleep(100 * time.Millisecond)
	cancel() // 通知 worker 退出
	time.Sleep(50 * time.Millisecond)

	// ---- 2. WithTimeout:超时自动取消 ----
	ctx2, cancel2 := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel2() // 建议 defer cancel 释放资源
	go worker(ctx2, "超时 worker")
	time.Sleep(200 * time.Millisecond)

	// ---- 3. WithDeadline:指定截止时间 ----
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()
	go worker(ctx3, "deadline worker")
	time.Sleep(150 * time.Millisecond)

	// ---- 4. WithValue:传递请求范围的值 ----
	type keyType string
	const traceIDKey keyType = "traceID"

	ctx4 := context.WithValue(context.Background(), traceIDKey, "req-abc-123")
	handleRequest(ctx4, traceIDKey)

	// ---- 5. 判断错误类型 ----
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel5()
	time.Sleep(20 * time.Millisecond)
	err := ctx5.Err()
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("原因:超时")
	}
	if errors.Is(err, context.Canceled) {
		fmt.Println("原因:被 cancel 调用")
	}
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] 退出: %v\n", name, ctx.Err())
			return
		default:
			// 干活
			time.Sleep(30 * time.Millisecond)
			fmt.Printf("[%s] tick\n", name)
		}
	}
}

func handleRequest(ctx context.Context, key any) {
	if v := ctx.Value(key); v != nil {
		fmt.Printf("处理请求, traceID=%v\n", v)
	}
}
