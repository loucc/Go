// Package main 演示接口。
//
// 学习要点:
//   - Go 的接口是**隐式实现**,不需要 implements 关键字
//   - 空接口 interface{} = any(Go 1.18+ 推荐 any)
//   - 类型断言 v, ok := i.(T);类型 switch
//   - 设计原则:接受接口,返回结构体;小接口优于大接口
//   - 接口底层是两个字:类型信息 + 数据指针
//   - nil 接口 vs 包含 nil 值的接口(著名陷阱)
//
// 运行:go run .
package main

import "fmt"

// 1. 定义接口(小而聚焦)
type Animal interface {
	Sound() string
	Name() string
}

// 2. 隐式实现(无需 implements)
type Dog struct{ name string }

func (d Dog) Sound() string { return "汪汪" }
func (d Dog) Name() string  { return d.name }

type Cat struct{ name string }

func (c Cat) Sound() string { return "喵" }
func (c Cat) Name() string  { return c.name }

// 3. 接受接口,返回结构体(经典设计模式)
func describe(a Animal) {
	fmt.Printf("%s 说:%s\n", a.Name(), a.Sound())
}

// 4. 类型断言的两种形式
func inspect(v any) {
	// 单返回值:失败会 panic
	// s := v.(string)

	// 双返回值:安全
	if s, ok := v.(string); ok {
		fmt.Printf("是字符串:%q\n", s)
		return
	}
	if n, ok := v.(int); ok {
		fmt.Printf("是整数:%d\n", n)
		return
	}
	fmt.Printf("未知类型 %T\n", v)
}

// 5. 类型 switch
func classify(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("bool: %v", x)
	case int, int32, int64:
		return fmt.Sprintf("整数: %v", x)
	case string:
		return fmt.Sprintf("字符串长度=%d", len(x))
	case Animal:
		return fmt.Sprintf("动物: %s", x.Name())
	default:
		return fmt.Sprintf("其他: %T", v)
	}
}

// 6. 演示 nil 接口的陷阱
type MyErr struct{ msg string }

func (e *MyErr) Error() string { return e.msg }

func mayFail() error {
	var p *MyErr = nil
	return p // 返回的接口"包裹了一个 nil 指针",但接口本身不是 nil!
}

func main() {
	// 隐式实现:任何有 Sound/Name 方法的类型都是 Animal
	describe(Dog{name: "旺财"})
	describe(Cat{name: "咪咪"})

	// 类型断言
	inspect("hello")
	inspect(42)
	inspect(3.14)

	// 类型 switch
	fmt.Println(classify(nil))
	fmt.Println(classify(true))
	fmt.Println(classify(100))
	fmt.Println(classify("golang"))
	fmt.Println(classify(Dog{name: "斑点"}))

	// nil 接口陷阱
	err := mayFail()
	if err != nil {
		// 会进入这里!因为接口内部有类型信息
		fmt.Println("⚠️ err != nil 但值是 nil:", err == nil, err)
	}
	// 正确写法:直接返回 nil,而不是返回一个 nil 指针
}
