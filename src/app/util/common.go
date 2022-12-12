package util

import (
	"os"
	"time"

	"github.com/revel/revel/logger"
)

var (
	Location  = "asia-northeast1"
	Zone      = "asia-northeast1-c"
	BuildDate string
	jst       *time.Location
)

const format = "2006 年 01 月 02 日 15:04:05"

func init() {
	logger.LogFunctionMap["stdoutjson"] =
		func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
			c.SetJson(os.Stdout, options)
		}
	logger.LogFunctionMap["stderrjson"] =
		func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
			c.SetJson(os.Stderr, options)
		}
	jst, _ = time.LoadLocation("Asia/Tokyo")
}

func DateToStr(value time.Time) string {
	return value.In(jst).Format(format)
}
