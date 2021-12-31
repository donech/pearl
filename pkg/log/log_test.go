package log

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	conf := Config{
		ServiceName: "log-test",
		Level:       "info",
		LevelColor:  true,
		Format:      "json",
		Stdout:      true,
		File: FileLogConfig{
			Filename:   "test.log",
			LogRotate:  true,
			MaxSize:    20,
			MaxAge:     20,
			MaxBackups: 10,
			BufSize:    20,
		},
		EncodeKey: EncodeKeyConfig{},
		SentryDSN: "",
	}
	_, err := New(conf)
	if err != nil {
		t.Error("创建 ginzap.logger 失败")
	}
	zap.S().Error(fmt.Sprint("Info log ", 2), zap.String("level", `{"a":"4","b":"5"}`))
}
