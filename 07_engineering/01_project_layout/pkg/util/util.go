// Package util 演示 pkg 目录中的公共工具。
package util

import "strings"

// Reverse 反转字符串(按 rune)。
func Reverse(s string) string {
	rs := []rune(s)
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	return string(rs)
}

// TitleWords 将 "hello world" 转为 "Hello World"(简单示例)。
func TitleWords(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
