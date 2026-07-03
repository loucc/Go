// Package main 演示嵌入(embedding)= 组合优于继承。
//
// 学习要点:
//   - Go 没有 extends/继承,只有组合
//   - 嵌入结构体 → 方法与字段被"提升"到外层
//   - 嵌入接口 → 组合出更大的接口
//   - 同名方法可以被外层覆盖(shadowing)
//
// 运行:go run .
package main

import "fmt"

// 1. 结构体嵌入
type Logger struct {
	prefix string
}

func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.prefix, msg)
}

// Service 嵌入了 Logger,自动拥有 Log 方法
type Service struct {
	*Logger // 嵌入指针(避免大对象复制)
	Name    string
}

// 2. 覆盖父方法(同名会遮蔽)
type LoudService struct {
	*Logger
	Name string
}

func (l *LoudService) Log(msg string) {
	// 显式调用嵌入的原始方法
	l.Logger.Log("LOUD: " + msg)
}

// 3. 接口嵌入(组合出大接口)
type Reader interface {
	Read(p []byte) (int, error)
}

type Writer interface {
	Write(p []byte) (int, error)
}

type Closer interface {
	Close() error
}

// 组合出 ReadWriteCloser(这就是标准库 io.ReadWriteCloser 的做法)
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// 4. 一个实现所有 3 个接口的类型
type File struct {
	name string
}

func (f *File) Read(p []byte) (int, error)  { fmt.Println("Read", f.name); return 0, nil }
func (f *File) Write(p []byte) (int, error) { fmt.Println("Write", f.name); return len(p), nil }
func (f *File) Close() error                { fmt.Println("Close", f.name); return nil }

func main() {
	// 嵌入:Service 直接拥有 Log 方法
	s := Service{
		Logger: &Logger{prefix: "svc"},
		Name:   "user-service",
	}
	s.Log("service started") // 相当于 s.Logger.Log(...)

	// 覆盖
	ls := LoudService{
		Logger: &Logger{prefix: "loud"},
		Name:   "loud-service",
	}
	ls.Log("hello")

	// 接口嵌入
	var rwc ReadWriteCloser = &File{name: "a.txt"}
	buf := []byte("hi")
	rwc.Write(buf)
	rwc.Read(buf)
	rwc.Close()
}
