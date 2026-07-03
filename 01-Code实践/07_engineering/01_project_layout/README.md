# 项目布局(Standard Go Project Layout)

参考 https://github.com/golang-standards/project-layout

```
project/
├── cmd/                     可执行程序入口(每个子目录是一个 main 包)
│   └── app/
│       └── main.go
├── internal/                私有代码(只有本模块能 import)
│   ├── service/             业务逻辑
│   ├── repository/          数据访问
│   └── handler/             HTTP/gRPC 处理器
├── pkg/                     可被外部项目引用的公共库(可选)
│   └── util/
├── api/                     API 定义(protobuf、OpenAPI)
├── configs/                 配置文件
├── deployments/             Dockerfile、k8s、docker-compose
├── scripts/                 构建/部署脚本
├── test/                    集成测试、e2e、测试数据
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 关键约定

- **`internal/`**:Go 编译器强制的可见性,里面的包只能被当前模块 import
- **`cmd/xxx/main.go`**:一个模块可以有多个可执行程序
- **`pkg/`**:如果不打算给别人 import,可以不用这个目录,直接把包放在项目根
- **`api/`**:proto 文件、swagger.json、gRPC 服务定义

## 分层示例

一个典型的 HTTP 服务分层:

```
handler (HTTP 层)   →   service (业务)   →   repository (数据)
    ↑ 依赖                  ↑ 依赖                ↑ 依赖
   API 结构体            domain 模型          domain 模型
```

依赖方向:上层依赖下层,下层不依赖上层;通过接口解耦。

## 本例文件

- `cmd/app/main.go` —— 程序入口,组装依赖
- `internal/service/user.go` —— 业务服务
- `pkg/util/util.go` —— 公共工具
