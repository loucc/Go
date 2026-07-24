// cmd/app/main.go —— 程序入口:组装依赖,启动。
//
// 注意 import 路径:
//
//	go-learning/07_engineering/01_project_layout/internal/service
//
// 因为项目根的 go.mod 定义了模块名 "go-learning"。
package main

import (
	"fmt"

	"go-learning/07_engineering/01_project_layout/internal/service"
	"go-learning/07_engineering/01_project_layout/pkg/util"
)

func main() {
	fmt.Println("=== Project Layout Demo ===")

	svc := service.New()
	svc.Create("Alice")
	svc.Create("Bob")
	svc.Create("Carol")

	for _, u := range svc.List() {
		fmt.Printf("  %d - %s (reversed: %s)\n",
			u.ID, u.Name, util.Reverse(u.Name))
	}

	u, err := svc.Get(2)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("Get(2) = %+v\n", u)
	}
}
