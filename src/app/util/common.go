package util

import (
	"time"
)

const format = "2006 年 01 月 02 日 15:04:05"

var jst *time.Location

func init() {
	jst, _ = time.LoadLocation("Asia/Tokyo")
}

func DateToStr(value time.Time) string {
	return value.In(jst).Format(format)
}
