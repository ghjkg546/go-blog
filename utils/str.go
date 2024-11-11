package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TruncateString(s string, maxRunes int) (string, bool) {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s, false
	}
	return string(runes[:maxRunes]), true
}

func formatTime(inputTime string) string {

	// 解析 ISO8601 格式时间字符串
	parsedTime, err := time.Parse(time.RFC3339Nano, inputTime)
	if err != nil {
		fmt.Println("解析时间错误:", err)
		return ""
	}

	// 将时间格式化为自定义的 "Y-m-d H:i:s" 格式
	formattedTime := parsedTime.Format("2006-01-02 15:04:05")

	fmt.Println("格式化后的时间:", formattedTime)
	return formattedTime
}

func TimestampToDateYmd(timestamp int64) string {

	// 解析 ISO8601 格式时间字符串
	t := time.Unix(timestamp, 0)

	// 格式化为日期字符串
	formattedTime := t.Format("2006-01-02 15:04:05")
	return formattedTime
}
