package logger

import (
	"fmt"
	"github.com/jimu-server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

var (
	Logger           *zap.Logger
	MultiWriteSyncer zapcore.WriteSyncer
)

func init() {
	FileName := "system.log"
	MaxSize := 100 // megabytes
	MaxBackups := 30
	MaxAge := 30 // days
	Level := "info"
	if config.Evn.App.Logger.Level != "" {
		Level = config.Evn.App.Logger.Level
	}
	var zapLevel zapcore.Level
	err := zapLevel.Set(Level)
	if err != nil {
		panic(err.Error())
		return
	}
	if config.Evn.App.Logger.FileName != "" {
		FileName = config.Evn.App.Logger.FileName
	}
	if config.Evn.App.Logger.MaxSize != 0 {
		MaxSize = config.Evn.App.Logger.MaxSize * 1024 * 1024 // bytes
	}
	if config.Evn.App.Logger.MaxBackups != 0 {
		MaxBackups = config.Evn.App.Logger.MaxBackups
	}
	if config.Evn.App.Logger.MaxAge != 0 {
		MaxAge = config.Evn.App.Logger.MaxAge
	}

	if strings.HasSuffix(FileName, ".log") {
		FileName = FileName[0 : len(FileName)-len(".log")]
	}
	// 创建控制台日志持久化
	consoleLog := &lumberjack.Logger{
		Filename:   FileName + ".log",
		MaxSize:    MaxSize, // megabytes
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge, //days
	}

	// 创建ERROR日志持久化
	errorLog := &lumberjack.Logger{
		Filename:   FileName + "-err.log",
		MaxSize:    MaxSize, // megabytes
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge, //days
	}
	// 创建持久化日志写入
	MultiWriteSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(consoleLog), zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(EncoderConfig(), MultiWriteSyncer, zapLevel)
	errCore := zapcore.NewCore(EncoderConfig(), zapcore.AddSync(errorLog), zapcore.ErrorLevel)
	Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}

func EncoderConfig() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "line",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,     // 日志换行符号
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// 自定义时间编码格式
			enc.AppendString(t.Format(time.DateTime))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
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
