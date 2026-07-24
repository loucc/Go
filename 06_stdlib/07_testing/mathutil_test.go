// mathutil_test.go 演示表驱动测试、子测试、Benchmark、Fuzz。
//
// 运行:
//
//	go test -v ./...
//	go test -bench=. ./...
//	go test -fuzz=FuzzAdd -fuzztime=5s ./...
package mathutil

import "testing"

// ---- 1. 基础单元测试 ----
func TestAdd_Basic(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add(2,3) = %d, want 5", got)
	}
}

// ---- 2. 表驱动测试(推荐)+ 子测试 ----
func TestAdd_Table(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"两正数", 1, 2, 3},
		{"含零", 0, 5, 5},
		{"负数", -1, -1, -2},
		{"混合", -3, 5, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// ---- 3. 用 Cleanup 释放资源 ----
func TestWithCleanup(t *testing.T) {
	t.Cleanup(func() {
		// 类似 defer,但在整个 t.Run 结束后调用
	})
}

// ---- 4. Benchmark ----
func BenchmarkFib(b *testing.B) {
	for b.Loop() { // Go 1.24+ 推荐写法
		_ = Fib(20)
	}
}

// ---- 5. Fuzz(Go 1.18+) ----
// go test -fuzz=FuzzAdd -fuzztime=5s
func FuzzAdd(f *testing.F) {
	f.Add(1, 2)
	f.Add(-5, 5)
	f.Fuzz(func(t *testing.T, a, b int) {
		sum := Add(a, b)
		if sum-b != a {
			t.Errorf("Add 反运算失败: (%d,%d)=%d", a, b, sum)
		}
	})
}
