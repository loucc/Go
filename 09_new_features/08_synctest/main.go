// Package main 演示 Go 1.25:testing/synctest。
//
// 痛点:测试里遇到 time.Sleep / time.After / context.WithTimeout 时,
// 要么真的睡等(测试变慢),要么改造代码注入时钟(侵入)。
//
// synctest 提供一个"虚拟时钟 + 全部 goroutine 阻塞时统一推进时间"的
// 沙箱环境,让涉及时间的并发逻辑可以确定性地、瞬时地跑完。
//
// 本文件是可运行的示例(非 _test.go),用于展示 API 形状;
// 真实使用请写在 _test.go 里,搭配 go test 运行。
//
//	func TestTimeout(t *testing.T) {
//	    synctest.Run(func() {
//	        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	        defer cancel()
//
//	        done := make(chan struct{})
//	        go func() {
//	            <-ctx.Done()  // 虚拟时间下 1s 一到立刻返回
//	            close(done)
//	        }()
//
//	        synctest.Wait()   // 等所有 goroutine 都阻塞,然后推进虚拟时间
//	        select {
//	        case <-done:
//	        default:
//	            t.Fatal("超时未触发")
//	        }
//	    })
//	}
package main

import "fmt"

func main() {
	fmt.Println("synctest 用于 _test.go,见本文件顶部注释里的示例。")
	fmt.Println("跑:go test ./... 时,涉及时间的测试会瞬时完成。")
}
