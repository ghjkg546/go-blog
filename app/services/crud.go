package services

import (
	"fmt"
	"time"
)

type crudService struct {
}

var CrudService = new(crudService)

func (mediaService *crudService) ParseStartEndTime(startTime string, endTime string) (int64, int64) {
	// 将日期字符串解析为 time.Time 对象
	layout := "2006-01-02" // 日期格式
	parsedTime, err := time.Parse(layout, startTime)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return 0, 0
	}

	parsedTime1, err := time.Parse(layout, endTime)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return 0, 0
	}

	// 转换为 10 位 Unix 时间戳
	startstamp := parsedTime.Unix()
	endstamp := parsedTime1.Unix() + (24*3600 - 1)

	return startstamp, endstamp
}
