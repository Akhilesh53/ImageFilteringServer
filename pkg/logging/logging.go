package logging

import (
	"os"
	"strings"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logInstance    *zap.Logger
	currWorkDir, _ = os.Getwd()
	logFileDir     = currWorkDir + "/logs/image_filter_server/"
	maxLogFileSize = 40
	maxLogFileAge  = 30
	maxLogFiles    = 100
	loggerOnce     sync.Once
)

func InitialiseLogger() *zap.Logger {
	loggerOnce.Do(func() {
		var (
			level       = getLevel(viper.GetString("LOG_LEVEL"))
			devConfig   = zap.NewDevelopmentEncoderConfig()
			prodConfig  = zap.NewProductionEncoderConfig()
			hostname, _ = os.Hostname()
		)

		devConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		prodConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

		if hostname != "" {
			hostname = "_" + hostname
		}

		filewriter := zapcore.AddSync(&FlushingLumberjackLogger{
			Logger: &lumberjack.Logger{
				Filename:   logFileDir + strings.ToLower(viper.GetString("PROCESS_NAME")) + ".log",
				MaxSize:    maxLogFileSize, // MB
				MaxAge:     maxLogFileAge,  // Days
				MaxBackups: maxLogFiles,    // No of files
				Compress:   false,          // disabled by default
			},
		})
		core := zapcore.NewCore(zapcore.NewJSONEncoder(prodConfig), filewriter, level)

		if viper.GetString("ENVIRONMENT") == "dev" {
			core = zapcore.NewTee(core, zapcore.NewCore(zapcore.NewConsoleEncoder(devConfig), zapcore.Lock(os.Stdout), level))
		}

		logInstance = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	})
	return logInstance
}

func Info(ctx *gin.Context, msg string, args ...zap.Field) {
	logInstance.Info(msg, append([]zap.Field{zap.String("uuid", ctx.GetString("uuid"))}, args...)...)
}

func Error(ctx *gin.Context, msg string, args ...zap.Field) {
	logInstance.Error(msg, append([]zap.Field{zap.String("uuid", ctx.GetString("uuid"))}, args...)...)
}

func Warn(ctx *gin.Context, msg string, args ...zap.Field) {
	logInstance.Warn(msg, append([]zap.Field{zap.String("uuid", ctx.GetString("uuid"))}, args...)...)
}

func Debug(ctx *gin.Context, msg string, args ...zap.Field) {
	logInstance.Debug(msg, append([]zap.Field{zap.String("uuid", ctx.GetString("uuid"))}, args...)...)
}

func Fatal(ctx *gin.Context, msg string, args ...zap.Field) {
	logInstance.Fatal(msg, append([]zap.Field{zap.String("uuid", ctx.GetString("uuid"))}, args...)...)
}

func getLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.ErrorLevel
	}
}
