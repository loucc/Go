// Package main 演示 go mod 与 go work 的核心命令。
//
// 本文件本身不运行复杂逻辑,主要作为命令速查的"可执行文档"。
// 实际操练请在终端执行下面的命令。
//
// ============================
// 一、go mod — 依赖管理
// ============================
//
// 1. 初始化模块(新项目)
//
//	go mod init github.com/yourname/project
//
// 2. 添加依赖(写进代码后自动拉取)
//
//	import "golang.org/x/sync/errgroup"
//	go mod tidy            // 自动添加缺失依赖 + 删除不再使用的
//
// 3. 查看当前依赖
//
//	go list -m all
//
// 4. 升级依赖
//
//	go get golang.org/x/sync@latest      // 升到最新
//	go get golang.org/x/sync@v0.7.0      // 升到指定版本
//
// 5. 降级 / 锁定版本
//
//	go get golang.org/x/sync@v0.5.0
//
// 6. 清理 & 校验
//
//	go mod tidy     // 清理无用依赖
//	go mod verify   // 校验 checksum
//
// 7. vendor(把依赖复制到项目内,CI 友好)
//
//	go mod vendor          // 复制依赖到 vendor/
//	go build -mod=vendor   // 用 vendor 目录编译
//
// ============================
// 二、go work — 多模块工作区(Go 1.18+)
// ============================
//
// 适用场景:本地同时开发多个相互依赖的模块,不想每次都 push 再 go get。
//
// 1. 初始化 workspace
//
//	go work init
//	go work use ./module-a
//	go work use ./module-b
//
// 生成的 go.work:
//
//	go 1.26
//
//	use (
//	    ./module-a
//	    ./module-b
//	)
//
// 2. 常用命令
//
//	go work use ./new-module    // 添加模块
//	go work edit -dropuse=./m   // 移除模块
//	go work sync                // 同步所有模块的 go.sum
//
// 3. 注意事项
//   - go.work 不应提交到 Git(个人开发环境)
//   - 如果要提交,团队约定一致即可
//   - GOWORK=off 可以临时禁用 workspace
//
// ============================
// 三、go.sum — 校验和
// ============================
//
// go.sum 记录每个依赖的 SHA-256 校验和,确保:
//   - 不同机器 / CI 拉到完全相同的代码
//   - 防篡改:如果远程代码被改,go build 会报错
//
// go.sum 应该提交到 Git!
//
// ============================
// 四、常见 go.mod 陷阱
// ============================
//
// 1. go 指令是"最低兼容版本",不是"目标版本"
//
//	go 1.22    // 意味着:本模块至少需要 Go 1.22
//
// 2. indirect 标记
//
//	require golang.org/x/text v0.14.0 // indirect
//
// // indirect 表示"本模块没有直接 import,但传递依赖需要"
// // go mod tidy 会自动管理,不要手动改
//
// 3. replace 指令(谨慎使用)
//
//	replace example.com/old => github.com/new v1.2.3
//
// 常见用途:fork 修复、本地开发调试
// 注意:replace 只在本模块生效,不会传递给依赖你的模块
//
// 运行:go run .
package main

import "fmt"

func main() {
	fmt.Println("=== Go Modules & Workspace 速查 ===")
	fmt.Println()
	fmt.Println("依赖管理:")
	fmt.Println("  go mod init <module-path>   初始化")
	fmt.Println("  go mod tidy                  清理 + 补全")
	fmt.Println("  go get <pkg>@latest          升级依赖")
	fmt.Println("  go list -m all               查看所有依赖")
	fmt.Println("  go mod verify                校验 checksum")
	fmt.Println("  go mod vendor                复制到 vendor/")
	fmt.Println()
	fmt.Println("多模块工作区:")
	fmt.Println("  go work init                 初始化 workspace")
	fmt.Println("  go work use ./module         添加模块")
	fmt.Println("  go work sync                 同步 go.sum")
	fmt.Println("  GOWORK=off go build          临时禁用")
	fmt.Println()
	fmt.Println("详见本文件源码注释。")
}
