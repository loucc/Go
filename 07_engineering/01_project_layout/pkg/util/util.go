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

// TitleWords 将 "hello world" 转为 "Hello World"。
// 按 rune 处理,支持 Unicode 字符。
func TitleWords(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		runes := []rune(w)
		if len(runes) > 0 {
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}
