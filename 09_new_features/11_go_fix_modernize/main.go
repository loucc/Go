// Package main 演示 Go 1.26:go fix 被重写成"现代化改写器"(modernizer)。
//
// 1.26 之前:go fix 只处理少数几个历史遗留问题,几乎没人用。
// 1.26 之后:基于 go vet 同一套 analyzer 框架,内置 20+ 修复器,
// 一条命令把你的代码升级到新版本习惯用法。
//
// 常用姿势:
//
//	# 看看会改哪些地方(不动源码)
//	go fix -diff ./...
//
//	# 直接应用修复(会改源码,建议先 commit)
//	go fix -fix ./...
//
//	# 只跑指定 analyzer,例如把 interface{} 换成 any
//	go fix -fix=modernize/anymap ./...
//
// 常见的现代化改写:
//   - interface{}                  →  any
//   - fmt.Sprintf 拼接              →  strings.Builder / strings.Join(能换的场合)
//   - for i := 0; i < n; i++       →  for range n            (1.22 整数 range)
//   - errors.New(fmt.Sprintf(...)) →  fmt.Errorf(...)
//   - sort.Slice                   →  slices.SortFunc         (泛型 slices 包)
//   - x = append(x[:i], x[i+1:]...) → slices.Delete           (泛型 slices 包)
//   - rand.Seed(...)               →  移除                    (1.20 起自动播种)
//   - ioutil.*                     →  io / os 对应函数
//   - `new(T); *p = v`             →  new(v)                 (1.26 new(表达式))
//
// 组合建议:
//
//	golangci-lint(风格 + bug)  +  go fix(半自动升级)  +  gopls(编辑器实时提示)
//
// 本文件不做具体演示 —— 到 01_basics ~ 06_stdlib 任意目录跑一下:
//
//	go fix -diff ./...
//
// 就能看到 modernizer 建议的改动。
package main

import "fmt"

func main() {
	fmt.Println("run: go fix -diff ./...  在项目根目录看会被现代化改写的地方。")
	fmt.Println("确认没问题再:go fix -fix ./...  自动应用。")
}
