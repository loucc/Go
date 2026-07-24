# Go 学习笔记

> Go 版本:1.26.4 | 模块:`go-learning` | 协议:MIT

一份结构化的 Go 语言学习资料,从基础语法到 1.26 新特性,每个知识点都配有可运行的代码示例。

## 快速开始

```bash
git clone https://github.com/loucc/Go.git
cd Go

# 运行任意示例
go run ./01_basics/01_hello
```

每个子目录都是独立的 `main` 包,也可以进入目录后执行 `go run .`。

## 项目结构

```
.
├── 01_basics/              基础语法
│   ├── 01_hello/           package / import / main / init
│   ├── 02_variables/       var / const / iota / 短声明
│   ├── 03_types/           基本类型、零值、类型转换
│   ├── 04_control_flow/    if / for / switch / goto
│   └── 05_builtin/         len / cap / make / new / min / max / clear
│
├── 02_types/               复合类型
│   ├── 01_arrays/          数组(值语义,长度是类型的一部分)
│   ├── 02_slices/          切片(三索引、扩容、slices 标准库)
│   ├── 03_maps/            map + maps 标准库
│   ├── 04_strings/         string / rune / strings.Builder / UTF-8
│   └── 05_pointers/        指针语义(无指针运算)
│
├── 03_functions/           函数、错误、泛型
│   ├── 01_functions/       多返回值、可变参数、闭包
│   ├── 02_defer/           LIFO、参数即时求值、性能开销
│   ├── 03_errors/          errors.Is / As / Unwrap、%w 包装
│   └── 04_generics/        类型参数、约束、~T
│
├── 04_oop/                 结构体、方法、接口
│   ├── 01_struct/          结构体、tag、匿名字段
│   ├── 02_methods/         值接收者 vs 指针接收者
│   ├── 03_interface/       隐式实现、类型断言、类型 switch
│   └── 04_embedding/       组合优于继承
│
├── 05_concurrency/         并发编程 ★
│   ├── 01_goroutine/       G-P-M 调度、启动开销
│   ├── 02_channel/         无缓冲 / 有缓冲、close 语义、单向
│   ├── 03_select/          多路复用、超时、非阻塞
│   ├── 04_sync/            Mutex / RWMutex / WaitGroup / Once / Pool
│   ├── 05_context/         WithCancel / WithTimeout / WithValue
│   └── 06_patterns/        worker pool、pipeline、fan-in / out
│
├── 06_stdlib/              标准库精讲
│   ├── 01_fmt/             格式化动词全集
│   ├── 02_io/              Reader / Writer 抽象
│   ├── 03_json/            encoding/json、tag、RawMessage
│   ├── 04_http/            net/http Server & Client(1.22+ 路由)
│   ├── 05_time/            time.Time / Duration / 时区
│   ├── 06_slog/            结构化日志(1.21+)
│   └── 07_testing/         单元测试、表驱动、基准测试、Fuzz
│
├── 07_engineering/         工程化实践
│   ├── 01_project_layout/  cmd / internal / pkg 约定
│   ├── 02_web_api/         REST API(net/http 1.22 路由)
│   └── 03_middleware/      中间件模式
│
├── 08_advanced/            进阶主题
│   ├── 01_reflect/         reflect.Value / reflect.Type
│   └── 02_unsafe/          unsafe.Pointer / Slice / String
│
└── 09_new_features/        Go 1.22–1.26 新特性
    ├── 01_range_int/       1.22  for range 整数
    ├── 02_loopvar/         1.22  循环变量新作用域
    ├── 03_http_routing/    1.22  net/http 增强路由
    ├── 04_range_func/      1.23  range over function(迭代器)
    ├── 05_iter_pkg/        1.23  iter 包
    ├── 06_generic_alias/   1.24  泛型类型别名
    ├── 07_waitgroup_go/    1.25  sync.WaitGroup.Go
    ├── 08_synctest/        1.25  testing/synctest
    ├── 09_new_with_expr/   1.26  new(表达式) 直接得到指针
    ├── 10_recursive_generics/ 1.26  自引用泛型(F-bounded)
    └── 11_go_fix_modernize/   1.26  go fix 现代化改写器
```

## 学习路线

| 阶段 | 模块 | 建议时长 |
|:---:|------|:---:|
| 1 | **基础语法** — 变量、类型、控制流、内置函数 | 1 周 |
| 2 | **复合类型** — 数组、切片、map、字符串、指针 | 1 周 |
| 3 | **函数与泛型** — 闭包、defer、错误处理、泛型 | 1 周 |
| 4 | **面向对象** — struct、方法、接口、嵌入 | 1 周 |
| 5 | **并发编程** ★ — goroutine、channel、select、sync、context | 1.5 周 |
| 6 | **标准库** — fmt、io、json、http、time、slog、testing | 2–3 周 |
| 7 | **工程化** — 项目布局、Web API、中间件 | 4–6 周 |
| 8 | **进阶** — reflect、unsafe | 按需 |
| 9 | **新特性** — Go 1.22 至 1.26 的语言演进 | 按需 |

> 核心内容(1–6)约 **2–3 个月**,达到中级工程能力约 **6 个月**。

## 配套练习建议

| 项目 | 涉及知识点 |
|------|-----------|
| CLI 工具(`wc` / `grep` / TODO) | 基础语法 + `flag` / `cobra` |
| HTTP JSON API | `net/http` + `encoding/json` + 中间件 |
| 并发爬虫 | goroutine + channel + `context` + rate limit |
| 短链服务 | `net/http` + Redis + Prometheus |
| 微型 RPC 框架 | `net` + `encoding` + `reflect` |

## 学习原则

- **接受接口,返回结构体**;小接口优于大接口
- **不要用 panic 做控制流**;错误应显式返回
- **Don't communicate by sharing memory; share memory by communicating**
- **接受 `context` 作为第一个参数**
- 优先**组合**而非继承
- 保持包名简短、包内符号简洁

## 环境要求

- Go 1.26.4+
- 编辑器:VSCode + Go 扩展 或 GoLand
- 推荐工具:`golangci-lint`、`delve`、`pprof`

## License

[MIT](LICENSE) © 2026 loucc
