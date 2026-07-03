// Package main 演示 Go 1.25:sync.WaitGroup.Go。
//
// 之前:
//
//	wg.Add(1)
//	go func() {
//	    defer wg.Done()
//	    doWork()
//	}()
//
// 1.25 起可以简写为:
//
//	wg.Go(doWork)
//
// 语法糖,但更不易忘记 Add/Done。
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		id := i
		wg.Go(func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("task %d done\n", id)
		})
	}

	wg.Wait()
	fmt.Println("all tasks done")
}
