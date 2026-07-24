// Package main 演示 log/slog(Go 1.21+ 官方结构化日志)。
//
// 学习要点:
//   - slog 支持结构化(key-value)日志
//   - 内置两种 handler:TextHandler / JSONHandler
//   - 通过 slog.With 创建带上下文的 logger
//   - 通过 slog.SetDefault 设置全局默认 logger
//
// 运行:go run .
package main

import (
	"errors"
	"log/slog"
	"os"
)

func main() {
	// ---- 1. TextHandler:人类可读 ----
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	textLogger.Info("server started", "port", 8080, "env", "prod")
	textLogger.Warn("slow query", "ms", 1500)

	// ---- 2. JSONHandler:适合 Loki / ELK 采集 ----
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	jsonLogger.Error("request failed",
		"err", errors.New("timeout"),
		"user_id", 42,
		"path", "/api/user",
	)

	// ---- 3. 带上下文的 logger ----
	reqLogger := jsonLogger.With(
		"request_id", "abc-123",
		"user", "alice",
	)
	reqLogger.Info("start")
	reqLogger.Info("done", "duration_ms", 42)

	// ---- 4. 分组 ----
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("db query",
		slog.Group("db",
			slog.String("host", "localhost"),
			slog.Int("port", 5432),
			slog.Duration("elapsed", 0),
		),
	)

	// ---- 5. 设置默认 logger ----
	slog.SetDefault(textLogger)
	slog.Info("这来自默认 logger")
}
