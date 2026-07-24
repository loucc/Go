// Package mathutil 提供一些示例函数用于测试演示。
package mathutil

// Add 两数之和。
func Add(a, b int) int { return a + b }

// Divide 除法(b=0 时返回 0, false)。
func Divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

// Fib 计算第 n 个斐波那契数。
func Fib(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
