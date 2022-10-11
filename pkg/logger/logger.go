package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gohub/pkg/app"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 处理日志相关逻辑

// Logger 全局Logger对象
var Logger *zap.Logger

// InitLogger日志初始化
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	// 获取日志写入介质
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)
	// 设置日志等级
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置错误，请修改config/log.go文件中log.level配置项")
	}
	// 初始化core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 初始化Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号，内部使用runtime.Caller
		zap.AddCallerSkip(1),              // 封装一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error才会显示stacktrace
	)

	// 将自定义的logger替换为全局的logger
	zap.ReplaceGlobals(Logger)
}

// getEncoder设置日志存储格式
func getEncoder() zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stackstrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加"\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写
		EncodeTime:     customTimeEncoder,              // 时间格式 2006-01-02 18:00:01
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间，单位为秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // caller短格式
	}
	// 本地环境配置
	if app.IsLocal() {
		// 标准输出关键字高亮
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		// 本地设置内置的Console解码器
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	// 线上环境使用JSON解码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 日志记录介质，使用标准输出+文件
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 配置按照日期记录日志文件
	if logType == "daily" {
		logname := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logname)
	}
	// 滚动日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	// 配置输出介质
	if app.IsLocal() {
		// 终端输出和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}
}
