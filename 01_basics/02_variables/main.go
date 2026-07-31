// Package main 演示变量与常量的多种声明方式。
//
// 学习要点:
//   - var:标准声明,可显式指定类型或让编译器推断 (var 变量名字 类型 = 表达式)
//   - :=  :短声明,仅在函数内部使用
//   - const + iota:枚举模式
//   - 类型别名 type A = B 与新类型 type A B 的区别
//   - Go 类型与 MySQL 字段的对应关系(GORM 实践)
//
// 运行:go run .
package main

import (
	"fmt"
	"strings"
	"time"
)

// 包级变量:必须用 var,不能用 :=
var (
	appName    = "go-learning"
	appVersion = "1.0.0"
	debug      bool // 零值 = false
)

// var 变量名字 类型 = 表达式
var age int8 = 20

// 常量:iota 从 0 开始,每行 +1
const (
	StatusPending  = iota // 0
	StatusRunning         // 1
	StatusFinished        // 2
	StatusFailed          // 3
)

// 位掩码模式(iota 结合位移)
const (
	FlagRead    = 1 << iota // 1
	FlagWrite               // 2
	FlagExecute             // 4
)

// 类型别名 vs 新类型
type UserID = int64 // 别名:UserID 和 int64 完全等价
type OrderID int64  // 新类型:OrderID 与 int64 不能隐式转换

func main() {
	// 1. 短声明(最常用,仅函数内)
	name := "Alice"
	age := 30

	// 2. 显式类型 (var 变量名字 类型 = 表达式)
	var height float64 = 1.75

	// 3. 多变量并行赋值(用于交换)
	a, b := 1, 2
	a, b = b, a

	// 4. 常量
	const Pi = 3.14159
	const MaxUsers int = 1000

	fmt.Printf("name=%s age=%d height=%.2f\n", name, age, height)
	fmt.Printf("swap: a=%d b=%d\n", a, b)
	fmt.Printf("Pi=%v MaxUsers=%v\n", Pi, MaxUsers)

	// iota 枚举
	fmt.Printf("Status: pending=%d running=%d finished=%d failed=%d\n",
		StatusPending, StatusRunning, StatusFinished, StatusFailed)

	// 位标志
	perm := FlagRead | FlagWrite
	fmt.Printf("perm=%b hasRead=%v hasExecute=%v\n",
		perm, perm&FlagRead != 0, perm&FlagExecute != 0)

	// 类型别名 vs 新类型
	var uid UserID = 100
	var oid OrderID = 200
	var i64 int64 = uid // ✅ 别名可以直接赋值
	// var i64_2 int64 = oid // ❌ 编译错误:cannot use oid (type OrderID) as type int64
	i64_2 := int64(oid) // 需要显式转换
	fmt.Printf("uid=%d i64=%d oid=%d i64_2=%d\n", uid, i64, oid, i64_2)

	_ = appName
	_ = appVersion
	_ = debug

	// ---- Go 类型 ↔ MySQL 字段映射(GORM 实践) ----
	demoMySQLTypeMapping()
}

// ---- 实战示例:用 GORM tag 标注 MySQL 列类型 ----
//
// 实际项目中 import "gorm.io/gorm" 后,
// AutoMigrate 会根据 struct tag 自动建表。
// 此处仅展示 tag 写法,不引入 GORM 依赖。
type Order struct {
	ID        int64   `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	UserID    int64   `gorm:"column:user_id;type:bigint;not null;index" json:"user_id"`
	OrderNo   string  `gorm:"column:order_no;type:varchar(64);uniqueIndex" json:"order_no"`
	Amount    string  `gorm:"column:amount;type:decimal(18,2);not null" json:"amount"` // 用 string 承接 decimal
	Exchange  string  `gorm:"column:exchange;type:decimal(18,6)" json:"exchange"`      // 外汇精度更高
	Status    int32   `gorm:"column:status;type:int;default:0" json:"status"`
	IsActive  bool    `gorm:"column:is_active;type:tinyint(1);default:1" json:"is_active"`
	Remark    string  `gorm:"column:remark;type:text" json:"remark"`
	Payload   []byte  `gorm:"column:payload;type:blob" json:"payload"`
	Counter   uint64  `gorm:"column:counter;type:bigint UNSIGNED" json:"counter"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;autoUpdateTime" json:"updated_at"`
}

func demoMySQLTypeMapping() {
	// 仅引用类型,避免 unused import
	_ = Order{}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 90))
	fmt.Println("Go 类型 ↔ MySQL 字段映射(GORM 实践)")
	fmt.Println(strings.Repeat("=", 90))

	rows := []struct {
		goType   string
		mysqlCol string
		gormTag  string
		scene    string
		note     string
	}{
		{"int64", "BIGINT(20)", `gorm:"type:bigint;primaryKey"`,
			"订单ID、用户ID、时间戳、金额(分)", "✅ 金融首选整型;对标 Java long"},
		{"uint64", "BIGINT UNSIGNED", `gorm:"type:bigint UNSIGNED"`,
			"计数器、序列号、权限掩码", "❌ 不要存金额(退款负数溢出)"},
		{"int32", "INT", `gorm:"type:int"`,
			"状态枚举、少量数字", "新项目尽量少用"},
		{"bool", "TINYINT(1)", `gorm:"type:tinyint(1)"`,
			"是否启用、有效标记", "true=1,false=0"},
		{"string", "VARCHAR(n)", `gorm:"type:varchar(64)"`,
			"订单号、手机号、账号名称", "长度按需设置"},
		{"string", "TEXT", `gorm:"type:text"`,
			"长文本、备注、原始报文", ""},
		{"[]byte", "BLOB", `gorm:"type:blob"`,
			"二进制、加密数据", ""},
		{"decimal.Decimal", "DECIMAL(M,N)", `gorm:"type:decimal(18,6)"`,
			"价格、费率、资产余额", "✅ 金融核心:DECIMAL(18,2)人民币;DECIMAL(18,6)外汇"},
		{"float64", "DOUBLE", `gorm:"type:double"`,
			"非金融统计计算", "❌ 严禁存储资金,精度丢失"},
		{"time.Time", "DATETIME", `gorm:"type:datetime;autoUpdateTime"`,
			"创建时间、更新时间", "导入 time 包"},
	}

	fmt.Printf("\n%-18s %-20s %-10s %s\n", "Go 类型", "MySQL 字段", "场景", "备注")
	fmt.Println(strings.Repeat("-", 90))
	for _, r := range rows {
		fmt.Printf("%-18s %-20s %-10s %s\n", r.goType, r.mysqlCol, r.scene, r.note)
	}

	fmt.Println()
	fmt.Println("GORM tag 示例:")
	for _, r := range rows {
		fmt.Printf("  %-18s → %s\n", r.goType, r.gormTag)
	}

	fmt.Println()
	fmt.Println("⚠️  关键原则:")
	fmt.Println("  1. 金额字段:优先 int64(单位:分),其次 decimal.Decimal;严禁 float64")
	fmt.Println("  2. ID 字段:统一 int64 + BIGINT,便于雪花算法/分布式 ID")
	fmt.Println("  3. 时间字段:用 time.Time + DATETIME;让 GORM 自动管理 CreatedAt/UpdatedAt")
	fmt.Println("  4. 布尔字段:MySQL 用 TINYINT(1),GORM 自动映射 bool")
	fmt.Println("  5. decimal 推荐库:github.com/shopspring/decimal")
}
