package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"ginDemo/config"
	"net/url"
	"sort"
	"strings"
	"time"
)

// 获取当前时间戳
func GetTimeUnix() int64 {
	return time.Now().Unix()
}

// MD5 方法
func MD5(str string) string {
	s := md5.New()
	if _, err := s.Write([]byte(str)); err != nil {
		return ""
	}
	return hex.EncodeToString(s.Sum(nil))
}

// 生成签名
func CreateSign(params url.Values) string {
	var key []string
	for k := range params {
		if k != "sn" && k != "ts" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	var builder strings.Builder
	for i, k := range key {
		if i > 0 {
			builder.WriteString("&")
		}
		fmt.Fprintf(&builder, "%v=%v", k, params.Get(k))
	}
	str := builder.String()

	// 自定义签名算法
	sign := MD5(MD5(str) + MD5(config.APP_NAME+config.APP_SECRET))
	return sign
}
