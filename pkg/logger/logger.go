package logger

import (
	"encoding/json"
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

// Dump调试专用，不会中断程序，会在终端输出warning消息
func Dump(value interface{}, msg ...string) {
	valueString := jsonString(value)
	// 判断第二个参数是否传参msg
	if len(msg) > 0 {
		Logger.Warn("Dump", zap.String(msg[0], valueString))
	} else {
		Logger.Warn("Dump", zap.String("data", valueString))
	}
}

// LogIf当err!=nil时记录error等级日志
func LogIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred:", zap.Error(err))
	}
}

// LogWarnIf 当err!=nil时记录warning日志
func LogWarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred:", zap.Error(err))
	}
}

// LogInfoIf 当err!=nil 时记录info等级日志
func LogInfoIf(err error) {
	if err != nil {
		Logger.Info("Error Occurred:", zap.Error(err))
	}
}

// Debug日志，详尽的程序日志
func Debug(moduleName string, fields ...zap.Field) {
	Logger.Debug(moduleName, fields...)
}

// Info日志
func Info(moduleName string, fields ...zap.Field) {
	Logger.Info(moduleName, fields...)
}

// Warn日志
func Warn(moduleName string, fields ...zap.Field) {
	Logger.Warn(moduleName, fields...)
}

// Error错误日志，不应该中断程序
func Error(moduleName string, fields ...zap.Field) {
	Logger.Error(moduleName, fields...)
}

// Fatal 级别和Error()一样，写完log后调用os.Exit(1)退出程序
func Fatal(moduleName string, fields ...zap.Field) {
	Logger.Fatal(moduleName, fields...)
}

// DebugString记录一条字符串类型的dubug日志
func DebugString(moduleName, name, msg string) {
	Logger.Debug(moduleName, zap.String(name, msg))
}
func InfoString(moduleName, name, msg string) {
	Logger.Info(moduleName, zap.String(name, msg))
}
func WarnString(moduleName, name, msg string) {
	Logger.Warn(moduleName, zap.String(name, msg))
}
func ErrorString(moduleName, name, msg string) {
	Logger.Error(moduleName, zap.String(name, msg))
}
func FatalString(moduleName, name, msg string) {
	Logger.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON记录对象类型的debug日志，使用json.Marshal进行编码
func DebugJSON(moduleName, name string, value interface{}) {
	Logger.Debug(moduleName, zap.String(name, jsonString(value)))
}
func InfoJSON(moduleName, name string, value interface{}) {
	Logger.Info(moduleName, zap.String(name, jsonString(value)))
}
func WarnJSON(moduleName, name string, value interface{}) {
	Logger.Warn(moduleName, zap.String(name, jsonString(value)))
}
func ErrorJSON(moduleName, name string, value interface{}) {
	Logger.Error(moduleName, zap.String(name, jsonString(value)))
}
func FatalJSON(moduleName, name string, value interface{}) {
	Logger.Fatal(moduleName, zap.String(name, jsonString(value)))
}

func jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		Logger.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
}

/*
Dump() —— 调试专用，会以结构化的形式输出到终端。且不会中断程序，使用 warn 等级（会有高亮）；
LogIf() / LogInfoIf() / LogWarnIf —— 减少我们代码中大量的 if err != nil { ... } 判断；
DebugString() 是语法糖，方便我们记录字符串类型的日志；
DebugJSON() 是语法糖，会使用 json.Marshal 记录的值，方便我们记录 struct 类型，以及其他类型的日志
*/
