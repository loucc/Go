// Package main 演示方法与值/指针接收者的选择。
//
// 学习要点:
//   - 方法就是带接收者的函数
//   - 值接收者:方法内修改不影响外部
//   - 指针接收者:方法内修改会反映到外部
//   - 一致性原则:一个类型的方法要么全值,要么全指针
//   - 选择指针接收者的时机:
//     ① 需要修改接收者;② 结构体较大(避免复制);③ 类型包含 sync.Mutex 等不可复制字段
//
// 运行:go run .
package main

import "fmt"

type Counter struct {
	count int
}

// 值接收者:内部修改是修改副本
func (c Counter) IncByValue() {
	c.count++
}

// 指针接收者:能真正修改
func (c *Counter) IncByPointer() {
	c.count++
}

// 输出用的方法通常用值接收者
func (c Counter) String() string {
	return fmt.Sprintf("Counter(%d)", c.count)
}

// 大结构体建议指针接收者
type LargeStruct struct {
	Data [1000]int
}

func (l *LargeStruct) DoSomething() { /* ... */ }

func main() {
	c := Counter{}
	c.IncByValue()
	c.IncByValue()
	c.IncByValue()
	fmt.Println("三次 IncByValue 后:", c) // 仍然是 0!

	c.IncByPointer()
	c.IncByPointer()
	c.IncByPointer()
	fmt.Println("三次 IncByPointer 后:", c) // 3

	// 指针也能调用值方法:自动解引用
	pc := &Counter{count: 10}
	fmt.Println("通过指针调用值方法:", pc.String())

	// 值能调用指针方法:自动取地址(前提是可寻址)
	c2 := Counter{}
	c2.IncByPointer() // 等价于 (&c2).IncByPointer()
	fmt.Println("c2:", c2)

	// ⚠️ 不可寻址的值不能调用指针方法
	// Counter{}.IncByPointer() // 编译错误
}
