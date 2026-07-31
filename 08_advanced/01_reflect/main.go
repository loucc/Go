// Package main 演示 reflect 反射。
//
// 学习要点:
//   - reflect.TypeOf / reflect.ValueOf
//   - Kind:reflect.Value 的分类(Int / String / Struct / Slice / ...)
//   - 修改值必须通过指针 + Elem
//   - 反射有性能代价,ORM/序列化等场景才用
//
// 运行:go run .
package main

import (
	"fmt"
	"reflect"
)

type Config struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"min=1,max=65535"`
	Debug    bool   `json:"debug"`
	optional string // 未导出字段,反射也能看见但不能设值
}

func main() {
	c := Config{Host: "localhost", Port: 8080, Debug: true}

	// ---- 1. 获取类型与值 ----
	t := reflect.TypeFor[Config]()
	v := reflect.ValueOf(c)
	fmt.Printf("Type: %s Kind: %s\n", t.Name(), t.Kind())

	// ---- 2. 遍历结构体字段 ----
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 未导出字段不能调用 Interface(),否则会 panic
		if !field.IsExported() {
			fmt.Printf("  %s (%s) = <unexported>  json=%q validate=%q\n",
				field.Name, field.Type,
				field.Tag.Get("json"), field.Tag.Get("validate"),
			)
			continue
		}

		fmt.Printf("  %s (%s) = %v  json=%q validate=%q\n",
			field.Name,
			field.Type,
			value.Interface(),
			field.Tag.Get("json"),
			field.Tag.Get("validate"),
		)
	}

	// ---- 3. 修改字段(必须传指针!) ----
	pv := reflect.ValueOf(&c).Elem() // 通过指针 + Elem 拿到可修改的 Value
	pv.FieldByName("Host").SetString("example.com")
	pv.FieldByName("Port").SetInt(9090)
	fmt.Printf("修改后: %+v\n", c)

	// ---- 4. 调用方法 ----
	m := reflect.ValueOf(&c).MethodByName("String")
	if m.IsValid() {
		result := m.Call(nil)
		fmt.Println("调用 String():", result[0].String())
	}

	// ---- 5. 动态创建切片 ----
	sliceType := reflect.SliceOf(reflect.TypeFor[int]())
	slice := reflect.MakeSlice(sliceType, 0, 5)
	for i := 1; i <= 3; i++ {
		slice = reflect.Append(slice, reflect.ValueOf(i*i))
	}
	fmt.Printf("动态切片: %v\n", slice.Interface())
}

// 给 Config 加一个方法用于反射调用演示
func (c Config) String() string {
	return fmt.Sprintf("Config{Host=%s Port=%d}", c.Host, c.Port)
}
