package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"ginDemo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 打印
func Print(i any) {
	fmt.Println("---")
	fmt.Println(i)
	fmt.Println("---")
}

// 返回JSON
func RetJson(code, msg string, data any, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}

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
		if k != "sn" {
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

// 验证签名
func VerifySign(c *gin.Context) {
	var method = c.Request.Method
	var ts int64
	var sn string
	var req url.Values

	if method == "GET" {
		req = c.Request.URL.Query()
		sn = c.Query("sn")
		var err error
		ts, err = strconv.ParseInt(c.Query("ts"), 10, 64)
		if err != nil {
			RetJson("500", "Invalid timestamp", "", c)
			return
		}

	} else if method == "POST" {
		c.Request.ParseForm()
		req = c.Request.PostForm
		sn = c.PostForm("sn")
		var err error
		ts, err = strconv.ParseInt(c.PostForm("ts"), 10, 64)
		if err != nil {
			RetJson("500", "Invalid timestamp", "", c)
			return
		}
	} else {
		RetJson("500", "Illegal requests", "", c)
		return
	}

	exp, err := strconv.ParseInt(config.API_EXPIRY, 10, 64)
	if err != nil {
		RetJson("500", "Config error", "", c)
		return
	}

	// 验证过期时间
	if ts > GetTimeUnix() || GetTimeUnix()-ts >= exp {
		RetJson("500", "Ts Error", "", c)
		return
	}

	// 验证签名
	if sn == "" || sn != CreateSign(req) {
		RetJson("500", "Sn Error", "", c)
		return
	}
}
