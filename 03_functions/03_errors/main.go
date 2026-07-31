// Package main 演示 Go 的错误处理。
//
// 学习要点:
//   - error 是接口:type error interface{ Error() string }
//   - 显式返回 + 逐层判断,不用异常
//   - errors.New / fmt.Errorf 创建 error
//   - fmt.Errorf 用 %w 包装 error,形成错误链
//   - errors.Is:判断链上是否包含某个 sentinel
//   - errors.As:提取链上的具体类型
//   - errors.Join:聚合多个错误(Go 1.20+)
//   - panic/recover 仅用于不可恢复的严重错误
//
// 运行:go run .
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// 1. 定义 sentinel error(哨兵错误,可比较)
var ErrNotFound = errors.New("not found")

// 2. 自定义 error 类型(携带更多上下文)
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s: %s", e.Field, e.Msg)
}

func findUser(id int) error {
	if id <= 0 {
		return &ValidationError{Field: "id", Msg: "must be positive"}
	}
	if id != 1 {
		// 用 %w 包装,保留原始错误
		return fmt.Errorf("user id=%d: %w", id, ErrNotFound)
	}
	return nil
}

func main() {
	// ---- 1. sentinel 错误比对 ----
	err := findUser(5)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("[Is] 用户不存在:", err)
	}

	// ---- 2. 提取自定义类型 ----
	err = findUser(-1)
	var verr *ValidationError
	if errors.As(err, &verr) {
		fmt.Printf("[As] 字段=%q 消息=%q\n", verr.Field, verr.Msg)
	}

	// ---- 3. 标准库错误链的实战:文件打开 ----
	_, err = os.Open("/nonexistent/path")
	if err != nil {
		fmt.Println("[os] raw error:", err)
		// os 的错误通常包裹了 fs.PathError
		var pe *fs.PathError
		if errors.As(err, &pe) {
			fmt.Printf("  op=%s path=%s underlying=%v\n", pe.Op, pe.Path, pe.Err)
		}
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("  (是 not-exist 类错误)")
		}
	}

	// ---- 4. 多重包装 ----
	e1 := errors.New("原始错误")
	e2 := fmt.Errorf("第二层: %w", e1)
	e3 := fmt.Errorf("第三层: %w", e2)
	fmt.Println("链:", e3)
	fmt.Println("Is e1?", errors.Is(e3, e1))

	// ---- 5. panic / recover(仅用于不可恢复错误) ----
	safeCall()

	fmt.Println("\n=== errors.Join(Go 1.20+) ===")
	demoJoin()
}

// ---- errors.Join:多错误聚合 ----
//
// 典型场景:并发任务中每个 worker 返回一个 error,
// 需要聚合后一起返回。Join 会跳过 nil,全部 nil 则返回 nil。
func demoJoin() {
	errs := []error{
		errors.New("连接超时"),
		nil, // Join 会跳过 nil
		errors.New("解析失败"),
		fmt.Errorf("字段 %q 无效: %w", "email", ErrNotFound),
	}

	combined := errors.Join(errs...)
	if combined != nil {
		fmt.Printf("聚合错误:\n%v\n", combined)
	}

	// Join 后的错误支持 errors.Is / errors.As 穿透
	if errors.Is(combined, ErrNotFound) {
		fmt.Println("→ 聚合错误中包含 ErrNotFound")
	}

	// 全部 nil 时 Join 返回 nil
	allNil := errors.Join(nil, nil, nil)
	fmt.Printf("全部 nil: %v\n", allNil)
}

func safeCall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recover: %v\n", r)
		}
	}()
	panic("something really bad")
}
