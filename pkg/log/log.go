package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//New new a zap.Logger
func New(conf Config) (*zap.Logger, error) {
	encoderConfig := getEncoderConfig(conf)
	if conf.ServiceName != "" {
		SystemName = conf.ServiceName
	}
	if conf.SystemTraceName != "" {
		SystemTraceName = conf.SystemTraceName
	}
	// 设置日志输出格式
	var encoder zapcore.Encoder
	switch conf.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	writeSyncer := make([]zapcore.WriteSyncer, 0, 2)
	writeSyncer = append(writeSyncer, zapcore.AddSync(os.Stderr))

	if conf.File.Filename != "" {
		// 添加日志切割归档功能
		hook := lumberjack.Logger{
			Filename:   conf.File.Filename,   // 日志文件路径
			MaxSize:    conf.File.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: conf.File.MaxBackups, // 日志文件最多保存多少个备份
			MaxAge:     conf.File.MaxAge,     // 文件最多保存多少天
			Compress:   false,                // 是否压缩
		}
		writeSyncer = append(writeSyncer, zapcore.AddSync(&hook))
	}

	core := zapcore.NewCore(
		encoder, // 编码器配置
		zapcore.NewMultiWriteSyncer(writeSyncer...),
		conf.level(), // 日志级别
	)

	// 构造日志
	logger := zap.New(core, zap.AddCaller(), zap.Development())
	logger = logger.Named(conf.ServiceName)

	// 将自定义的logger替换为全局的logger
	zap.ReplaceGlobals(logger)
	return logger, nil
}

func getEncoderConfig(conf Config) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	if conf.EncodeKey.TimeKey != "" {
		encoderConfig.TimeKey = conf.EncodeKey.TimeKey
	}
	if conf.EncodeKey.NameKey != "" {
		encoderConfig.NameKey = conf.EncodeKey.NameKey
	}
	if conf.EncodeKey.LevelKey != "" {
		encoderConfig.LevelKey = conf.EncodeKey.LevelKey
	}
	if conf.EncodeKey.CallerKey != "" {
		encoderConfig.CallerKey = conf.EncodeKey.CallerKey
	}
	if conf.EncodeKey.MessageKey != "" {
		encoderConfig.MessageKey = conf.EncodeKey.MessageKey
	}
	if conf.EncodeKey.StacktraceKey != "" {
		encoderConfig.StacktraceKey = conf.EncodeKey.StacktraceKey
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	return encoderConfig
}
