package util

import (
	"time"
)

var (
	Date string
	jst  *time.Location
)

const format = "2006 年 01 月 02 日 15:04:05"

func init() {
	jst, _ = time.LoadLocation("Asia/Tokyo")
}

func DateToStr(value time.Time) string {
	return value.In(jst).Format(format)
}
