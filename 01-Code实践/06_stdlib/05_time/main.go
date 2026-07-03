// Package main 演示 time 包。
//
// 学习要点:
//   - time.Time 是值类型
//   - time.Duration 底层是 int64 纳秒
//   - 格式化用的参考时间:2006-01-02 15:04:05 (记法:1 2 3 4 5 6 7)
//   - 时区处理:time.UTC / time.Local / LoadLocation
//   - Timer 与 Ticker
//
// 运行:go run .
package main

import (
	"fmt"
	"time"
)

func main() {
	// ---- 1. 当前时间 ----
	now := time.Now()
	fmt.Println("now:", now)
	fmt.Println("Unix ts:", now.Unix())
	fmt.Println("UnixMilli:", now.UnixMilli())

	// ---- 2. 构造 time.Time ----
	t := time.Date(2026, 7, 2, 10, 30, 0, 0, time.Local)
	fmt.Println("t:", t)

	// ---- 3. 格式化(参考时间要记牢!) ----
	// Go 用具体的时间点做占位:2006-01-02 15:04:05
	fmt.Println("Format YYYY-MM-DD:", now.Format("2006-01-02"))
	fmt.Println("Format 时分秒 :", now.Format("15:04:05"))
	fmt.Println("Format 完整   :", now.Format("2006-01-02 15:04:05.000"))
	fmt.Println("RFC3339      :", now.Format(time.RFC3339))

	// ---- 4. 解析 ----
	parsed, err := time.Parse("2006-01-02", "2026-07-02")
	if err == nil {
		fmt.Println("Parse:", parsed)
	}

	// ---- 5. Duration 与算术 ----
	d := 2*time.Hour + 30*time.Minute
	fmt.Println("Duration:", d, "秒数:", d.Seconds())

	future := now.Add(d)
	fmt.Println("2.5 小时后:", future.Format(time.RFC3339))

	diff := future.Sub(now)
	fmt.Println("Sub 差值:", diff)

	// ---- 6. 时区 ----
	loc, _ := time.LoadLocation("America/New_York")
	fmt.Println("北京时间:", now.In(time.Local).Format("15:04"))
	fmt.Println("纽约时间:", now.In(loc).Format("15:04"))
	fmt.Println("UTC   :", now.UTC().Format("15:04"))

	// ---- 7. Timer & Ticker ----
	// Timer:一次性
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	fmt.Println("timer 触发")

	// Ticker:周期性(用完必须 Stop!)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	count := 0
	for range ticker.C {
		count++
		fmt.Println("tick", count)
		if count >= 3 {
			break
		}
	}
}
