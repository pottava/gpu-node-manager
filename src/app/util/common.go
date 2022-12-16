package util

import (
	"context"
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

func ProjectID() string {
	if candidate, found := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); found {
		return candidate
	}
	meta := InstanceMetadata(context.Background())
	if value, ok := meta["project_id"]; ok {
		return value
	}
	panic("project id was not found")
}

func RunRevision() string {
	if candidate, found := os.LookupEnv("K_REVISION"); found {
		return candidate
	}
	return "local"
}

func AppStage() string {
	if candidate, found := os.LookupEnv("STAGE"); found {
		return candidate
	}
	return "local"
}

func DateToStr(value time.Time) string {
	return value.In(jst).Format(format)
}
