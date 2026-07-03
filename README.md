# Go 1.26.4 学习计划

> 参考项目:https://github.com/xinliangnote/Go
> Go 版本:1.26.4

## 目录结构

```
.
├── 01_basics/          基础语法(1 周)
├── 02_types/           复合类型(1 周)
├── 03_functions/       函数、错误、泛型(1 周)
├── 04_oop/             结构体、方法、接口(1 周)
├── 05_concurrency/     并发编程(1.5 周)★重头戏
├── 06_stdlib/          标准库精讲(2-3 周)
├── 07_engineering/     工程化实践(4-6 周)
├── 08_advanced/        进阶主题(按需)
└── 09_new_features/    Go 1.22-1.26 新特性
```

## 运行方式

每个子目录都是独立的 `main` 包,进入目录后执行:

```bash
cd 01_basics/01_hello
go run .
```

或者从项目根目录运行:

```bash
go run ./01_basics/01_hello
```

---

## 学习路线(9 阶段)

### 阶段 0:环境与工具链
- [x] 安装 Go 1.26.3
- [ ] 掌握 `go mod`、`go build`、`go run`、`go test`、`go vet`、`go fmt`
- [ ] 熟悉 `golangci-lint`、`delve`、`pprof`
- [ ] 配置 IDE(VSCode + Go 扩展,或 GoLand)

### 阶段 1:基础语法 [01_basics/]
- 01_hello       —— package/import/main
- 02_variables   —— var / const / iota / 短声明
- 03_types       —— 基本类型、零值、类型转换
- 04_control_flow —— if / for / switch / goto
- 05_builtin     —— len / cap / make / new / min / max / clear

### 阶段 2:复合类型 [02_types/]
- 01_arrays   —— 数组(值语义,长度是类型的一部分)
- 02_slices   —— 切片(重点!三索引、扩容、slices 标准库)
- 03_maps     —— map + maps 标准库
- 04_strings  —— string、rune、strings.Builder、UTF-8
- 05_pointers —— 指针语义(无指针运算)

### 阶段 3:函数、错误、泛型 [03_functions/]
- 01_functions —— 多返回值、可变参数、闭包
- 02_defer     —— LIFO、参数即时求值、性能开销
- 03_errors    —— errors.Is/As/Unwrap、%w 包装
- 04_generics  —— 类型参数、约束、~T

### 阶段 4:面向"接口"编程 [04_oop/]
- 01_struct    —— 结构体、tag、匿名字段
- 02_methods   —— 值接收者 vs 指针接收者
- 03_interface —— 隐式实现、类型断言、类型 switch
- 04_embedding —— 组合优于继承

### 阶段 5:并发编程 [05_concurrency/]  ★最重要
- 01_goroutine —— G-P-M 调度、启动开销
- 02_channel   —— 无缓冲/有缓冲、close 语义、单向
- 03_select    —— 多路复用、超时、非阻塞
- 04_sync      —— Mutex/RWMutex/WaitGroup/Once/Pool
- 05_context   —— WithCancel/WithTimeout/WithValue
- 06_patterns  —— worker pool、pipeline、fan-in/out

### 阶段 6:标准库 [06_stdlib/]
- 01_fmt      —— 格式化动词全集
- 02_io       —— Reader/Writer 抽象
- 03_json     —— encoding/json、tag、RawMessage
- 04_http     —— net/http Server & Client(1.22+ 路由)
- 05_time     —— time.Time / Duration / 时区
- 06_slog     —— 结构化日志(1.21+)
- 07_testing  —— 单元测试、表驱动、基准测试、Fuzz

### 阶段 7:工程化实践 [07_engineering/]
- 01_project_layout —— cmd / internal / pkg 约定
- 02_web_api        —— REST API(net/http 1.22 路由)
- 03_middleware     —— 中间件模式

### 阶段 8:进阶主题 [08_advanced/]
- 01_reflect —— reflect.Value / reflect.Type
- 02_unsafe  —— unsafe.Pointer、Slice、String

### 阶段 9:Go 1.22-1.26 新特性 [09_new_features/]
- 01_range_int         —— 1.22:for range 整数
- 02_loopvar           —— 1.22:循环变量新作用域
- 03_http_routing      —— 1.22:net/http 增强路由
- 04_range_func        —— 1.23:range over function(迭代器)
- 05_iter_pkg          —— 1.23:iter 包
- 06_generic_alias     —— 1.24:泛型类型别名
- 07_waitgroup_go      —— 1.25:sync.WaitGroup.Go
- 08_synctest          —— 1.25:testing/synctest
- 09_new_with_expr     —— 1.26:new(表达式) 直接得到指向该值的指针
- 10_recursive_generics —— 1.26:自引用泛型(F-bounded)
- 11_go_fix_modernize  —— 1.26:go fix 现代化改写(20+ analyzer)

---

## 时间预估(全职学习)

| 阶段 | 时间 |
|------|------|
| 0 环境 + 1 基础 + 2 复合类型 | 2-3 周 |
| 3 函数 + 4 OOP | 2 周 |
| 5 并发 ★ | 1.5 周 |
| 6 标准库(边学边做) | 2-3 周 |
| 7 工程化 | 4-6 周 |
| 8-9 进阶 + 新特性 | 按需 |

**总计**:掌握核心约 3 个月;达到中级工程能力约 6 个月。

---

## 推荐配套项目(照 xinliangnote/Go 走)

1. **CLI 工具**:实现 `wc` / `grep` / TODO CLI(基础 + flag/cobra)
2. **HTTP JSON API**:用户注册登录 + JWT + MySQL
3. **并发爬虫**:worker pool + rate limit + context 取消
4. **短链服务**:Redis + Gin + Prometheus
5. **微型 RPC 框架**:自己撸一个 gRPC-like

---

## 关键学习原则

- **接受接口,返回结构体**;小接口优于大接口
- **不要用 panic 做控制流**;错误应显式返回
- **Don't communicate by sharing memory; share memory by communicating**(用 channel 而不是共享内存)
- **接受 context 作为第一个参数**
- 优先 **组合**,而非继承
- 保持包名简短、包内符号简洁
