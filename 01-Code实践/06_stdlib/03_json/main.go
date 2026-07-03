// Package main 演示 encoding/json。
//
// 学习要点:
//   - json.Marshal / Unmarshal 是序列化/反序列化入口
//   - tag:字段名映射、omitempty、"-" 忽略
//   - 未知字段默认忽略;可以用 DisallowUnknownFields 严格模式
//   - 大小写:仅**导出字段**(首字母大写)才会被序列化
//   - json.RawMessage 延迟解析
//   - json.Number 用于不损失精度的数字
//
// 运行:go run .
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Address struct {
	City   string `json:"city"`
	Street string `json:"street,omitempty"` // 空值省略
}

type User struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"-"` // 永不序列化
	Emails   []string `json:"emails,omitempty"`
	Address  Address  `json:"address"`
}

func main() {
	// ---- 1. Marshal ----
	u := User{
		ID:      1,
		Name:    "Alice",
		Emails:  []string{"a@x.com", "a@y.com"},
		Address: Address{City: "Beijing"},
	}
	b, _ := json.Marshal(u)
	fmt.Println("紧凑:", string(b))

	// 带缩进
	b, _ = json.MarshalIndent(u, "", "  ")
	fmt.Println("缩进:")
	fmt.Println(string(b))

	// ---- 2. Unmarshal ----
	jsonStr := `{"id":2,"name":"Bob","address":{"city":"Shanghai","street":"南京路"}}`
	var u2 User
	if err := json.Unmarshal([]byte(jsonStr), &u2); err != nil {
		panic(err)
	}
	fmt.Printf("解析: %+v\n", u2)

	// ---- 3. 未知字段严格模式 ----
	dec := json.NewDecoder(strings.NewReader(`{"id":3,"unknown_field":true}`))
	dec.DisallowUnknownFields()
	var u3 User
	if err := dec.Decode(&u3); err != nil {
		fmt.Println("严格模式错误:", err)
	}

	// ---- 4. 大数字精度:json.Number ----
	dec2 := json.NewDecoder(strings.NewReader(`{"big":10000000000000000123}`))
	dec2.UseNumber()
	var m map[string]any
	dec2.Decode(&m)
	fmt.Printf("json.Number 精度:%v (类型 %T)\n", m["big"], m["big"])

	// ---- 5. RawMessage:延迟解析 ----
	type Envelope struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	raw := `{"type":"user","payload":{"id":10,"name":"Carol"}}`
	var env Envelope
	json.Unmarshal([]byte(raw), &env)
	fmt.Println("Envelope.Type =", env.Type)
	// 根据 Type 决定 Payload 具体解析成什么
	if env.Type == "user" {
		var payload User
		json.Unmarshal(env.Payload, &payload)
		fmt.Printf("Payload = %+v\n", payload)
	}

	// ---- 6. 流式输出 ----
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.Encode(u)
	fmt.Println("Encoder:", buf.String())
}
