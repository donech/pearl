package log

import "go.uber.org/zap/zapcore"

type Config struct {
	// 业务服务名称，如果多个业务日志在ELK中聚合，此字段就有用了
	ServiceName string `yaml:"serviceName"`
	// 日志级别.
	Level string `yaml:"level"`
	// 日志级别字段开启颜色功能
	LevelColor bool `yaml:"levelColor"`
	// Log format. one of json or plain.
	Format string `yaml:"format"`
	// 是否输出到控制台.
	Stdout bool `yaml:"stdout"`
	// File log config.
	File FileLogConfig `yaml:"file"`
	//EncodeKey EncodeKeys
	EncodeKey EncodeKeyConfig `yaml:"encodeKey"`
	//SentryDSN Sentry 的 DSN地址，如果配置了此参数，warn 级别以上的错误会发送sentry
	SentryDSN string `yaml:"sentryDSN"`
	//SystemName 系统名称
	SystemName string `yaml:"systemName"`
	//SystemTraceName 系统日志（非业务相关） traceID 对应的值
	SystemTraceName string `yaml:"systemTraceName"`
}

//level 获取日志级别，默认是Info
func (c *Config) level() zapcore.Level {
	level := zapcore.InfoLevel
	if c.Level == "" {
		return level
	}
	if err := level.Set(c.Level); err != nil {
		panic(err)
	}
	return level
}

//FileLogConfig serializes file log related config.
type FileLogConfig struct {
	// Filename 日志文件路径.
	Filename string `yaml:"filename"`
	// LogRotate Is log rotate enabled.
	LogRotate bool `yaml:"logRotate"`
	// MaxSize size for a single file, in MB.
	MaxSize int `yaml:"maxSize"`
	// MaxAge log keep days, default is never deleting.
	MaxAge int `yaml:"maxAge"`
	// MaxBackups  number of old log files to retain.
	MaxBackups int `yaml:"maxBackups"`
	// BufSize  size of bufio.Writer
	BufSize int `yaml:"bufSize"`
}

type EncodeKeyConfig struct {
	// TimeKey format.
	TimeKey string `yaml:"timeKey"`
	// LevelKey format.
	LevelKey string `yaml:"levelKey"`
	// NameKey format.
	NameKey string `yaml:"nameKey"`
	// CallerKey format.
	CallerKey string `yaml:"callerKey"`
	// MessageKey format.
	MessageKey string `yaml:"messageKey"`
	//StacktraceKey StacktraceKey
	StacktraceKey string `yaml:"stacktraceKey"`
}
