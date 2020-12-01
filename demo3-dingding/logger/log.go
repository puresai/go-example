package logger

import (
	"io"
	"log"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
    "github.com/spf13/viper"
)

var Logger *zap.Logger 
var LogLevel string
var FileFormat string

// 初始化日志 logger
func init() {
	// 设置一些基本日志格式
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	encoder := zapcore.NewConsoleEncoder(config)

	FileFormat, saveType, LogLevel := "%Y%m%d", "one", "info"

	if viper.IsSet("log.file_format") {
		FileFormat = viper.GetString("log.file_format")
	}

	if viper.IsSet("log.level") {
		LogLevel = viper.GetString("log.level")
	}

	if viper.IsSet("log.save_type") {
		saveType = viper.GetString("log.save_type")
	}

	logLevel := zap.DebugLevel
	switch LogLevel {
		case "debug":
			logLevel = zap.DebugLevel
		case "info":
			logLevel = zap.InfoLevel
		case "error":
			logLevel = zap.ErrorLevel
		default:
			logLevel = zap.InfoLevel
	}

	switch saveType {
		case "level":
			Logger = getLevelLogger(encoder, logLevel, FileFormat)
		default:
			Logger = getOneLogger(encoder, logLevel, FileFormat)
	}

	
}

func getLevelLogger(encoder zapcore.Encoder, logLevel zapcore.Level, fileFormat string) *zap.Logger {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && lvl >= logLevel
	})

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && lvl >= logLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= logLevel
	})
	// 获取 info、warn日志文件的io.Writer 抽象 getLoggerWriter() 在下方实现
	infoWriter := getLoggerWriter("./log/info", fileFormat)
	errorWriter := getLoggerWriter("./log/error", fileFormat)
	debugWriter := getLoggerWriter("./log/debug", fileFormat)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func getOneLogger(encoder zapcore.Encoder, logLevel zapcore.Level, fileFormat string) *zap.Logger {
	infoWriter := getLoggerWriter("./log/info", fileFormat)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && lvl >= logLevel
	})

	core := zapcore.NewTee(
		// 将info及以下写入logPath,  warn及以上写入errPath
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func getLoggerWriter(filename, fileFormat string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 file_YYmmddHH.log
	hook, err := rotatelogs.New(
		filename+fileFormat+".log",
		rotatelogs.WithLinkName(filename),
		// 保存天数
		rotatelogs.WithMaxAge(time.Hour*24*30),
		// 切割频率24小时
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}

// logs.Debug(...)
func Debug(format string, v ...interface{}) {
	Logger.Sugar().Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	Logger.Sugar().Infof(format, v...)
}

func Error(format string, v ...interface{}) {
	Logger.Sugar().Errorf(format, v...)
}