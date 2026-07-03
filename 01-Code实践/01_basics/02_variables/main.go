// Package main 演示变量与常量的多种声明方式。
//
// 学习要点:
//   - var:标准声明,可显式指定类型或让编译器推断
//   - :=  :短声明,仅在函数内部使用
//   - const + iota:枚举模式
//   - 类型别名 type A = B 与新类型 type A B 的区别
//
// 运行:go run .
package main

import "fmt"

// 包级变量:必须用 var,不能用 :=
var (
	appName    = "go-learning"
	appVersion = "1.0.0"
	debug      bool // 零值 = false
)

// 常量:iota 从 0 开始,每行 +1
const (
	StatusPending  = iota // 0
	StatusRunning         // 1
	StatusFinished        // 2
	StatusFailed          // 3
)

// 位掩码模式(iota 结合位移)
const (
	FlagRead    = 1 << iota // 1
	FlagWrite               // 2
	FlagExecute             // 4
)

// 类型别名 vs 新类型
type UserID = int64 // 别名:UserID 和 int64 完全等价
type OrderID int64  // 新类型:OrderID 与 int64 不能隐式转换

func main() {
	// 1. 短声明(最常用,仅函数内)
	name := "Alice"
	age := 30

	// 2. 显式类型
	var height float64 = 1.75

	// 3. 多变量并行赋值(用于交换)
	a, b := 1, 2
	a, b = b, a

	// 4. 常量
	const Pi = 3.14159
	const MaxUsers int = 1000

	fmt.Printf("name=%s age=%d height=%.2f\n", name, age, height)
	fmt.Printf("swap: a=%d b=%d\n", a, b)
	fmt.Printf("Pi=%v MaxUsers=%v\n", Pi, MaxUsers)

	// iota 枚举
	fmt.Printf("Status: pending=%d running=%d finished=%d failed=%d\n",
		StatusPending, StatusRunning, StatusFinished, StatusFailed)

	// 位标志
	perm := FlagRead | FlagWrite
	fmt.Printf("perm=%b hasRead=%v hasExecute=%v\n",
		perm, perm&FlagRead != 0, perm&FlagExecute != 0)

	// 类型别名 vs 新类型
	var uid UserID = 100
	var oid OrderID = 200
	var i64 int64 = uid // ✅ 别名可以直接赋值
	// var i64_2 int64 = oid // ❌ 编译错误:cannot use oid (type OrderID) as type int64
	i64_2 := int64(oid) // 需要显式转换
	fmt.Printf("uid=%d i64=%d oid=%d i64_2=%d\n", uid, i64, oid, i64_2)

	_ = appName
	_ = appVersion
	_ = debug
}
