package logger

import (
	"fmt"
	"github.com/jimu-server/config/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	Logger *zap.Logger
)

func init() {
	if config.Evn.App.Logger.Level == "" {
		zap.L().Panic("aurora. zap.level error")
	}
	var zapLevel zapcore.Level
	err := zapLevel.Set(config.Evn.App.Logger.Level)
	if err != nil {
		zap.L().Panic("aurora. zap.level error")
		return
	}

	if config.Evn.App.Logger.FileName == "" {
		zap.L().Panic("app.zap.filename error")
	}
	if config.Evn.App.Logger.MaxSize == 0 {
		zap.L().Panic("app.zap.maxsize error")
	}
	if config.Evn.App.Logger.MaxBackups == 0 {
		zap.L().Panic("app.zap.maxage error")
	}
	if config.Evn.App.Logger.MaxAge == 0 {
		zap.L().Panic("app.zap.maxbackups error")
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Evn.App.Logger.FileName,
		MaxSize:    config.Evn.App.Logger.MaxSize, // megabytes
		MaxBackups: config.Evn.App.Logger.MaxBackups,
		MaxAge:     config.Evn.App.Logger.MaxAge, //days
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(encoderConfig(), writeSyncer, zapLevel)
	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}

func encoderConfig() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "line",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,        // 日志换行符号
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// 自定义时间编码格式
			enc.AppendString(t.Format(time.DateTime))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewConsoleEncoder(config)
}

func Info(format string, a ...any) {
	Logger.Info(fmt.Sprintf(format, a...))
}

func Wring(format string, a ...any) {
	Logger.Warn(fmt.Sprintf(format, a...))
}
func Debug(format string, a ...any) {
	Logger.Debug(fmt.Sprintf(format, a...))
}
func Error(format string, a ...any) {
	Logger.Error(fmt.Sprintf(format, a...))
}