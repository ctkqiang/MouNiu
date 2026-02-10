package utilities

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// LogLevel 定义了日志的严重程度级别
type LogLevel int

const (
	APP_NAME = "某牛"
	VERSION  = "1.0.0"
)

// 日志级别枚举
const (
	DEBUG    LogLevel = iota // 调试级：用于开发阶段的详细信息
	INFO                     // 信息级：常规业务运行信息
	WARN                     // 警告级：可能存在潜在问题的异常情况
	ERROR                    // 错误级：影响业务逻辑的严重问题
	VVERBOSE                 // 冗余级：极其详细的底层追踪信息
)

// ANSI 颜色转义字符，用于在终端输出彩色日志
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPink   = "\033[35m"
)

var (
	// CurrentLevel 记录当前系统运行的日志级别，默认为 INFO
	CurrentLevel = INFO
)

// SetLogLevel 根据传入的字符串动态设置日志级别
func SetLogLevel(levelStr string) {
	level := strings.ToUpper(levelStr)
	switch level {
	case "DEBUG":
		CurrentLevel = DEBUG
	case "INFO":
		CurrentLevel = INFO
	case "WARN":
		CurrentLevel = WARN
	case "ERROR":
		CurrentLevel = ERROR
	case "VVERBOSE":
		CurrentLevel = VVERBOSE
	}
}

func init() {
	// 从环境变量 LOG_LEVEL 中初始化日志级别
	SetLogLevel(os.Getenv("LOG_LEVEL"))
}

// Log 是核心日志输出函数，负责格式化输出内容并处理日志级别过滤
func Log(level LogLevel, format string, a ...interface{}) {
	// 过滤掉低于当前设置级别的日志
	if level < CurrentLevel {
		return
	}

	levelStr := ""
	color := ""

	// 根据级别分配标签和颜色
	switch level {
	case DEBUG:
		levelStr = "调试"
		color = colorYellow
	case INFO:
		levelStr = "信息"
		color = colorBlue
	case WARN:
		levelStr = "警告"
		color = colorPink
	case ERROR:
		levelStr = "错误"
		color = colorRed
	case VVERBOSE:
		levelStr = "冗长"
		color = "" // 默认颜色
	}

	// 获取当前时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, a...)

	// 输出带有应用名称和时间戳的格式化日志
	if color != "" {
		fmt.Printf("[GIN-debug] %s[%s] [%s] %s%s\n", color, timestamp, levelStr, msg, colorReset)
	} else {
		fmt.Printf("[GIN-debug] [%s] [%s] %s\n", timestamp, levelStr, msg)
	}
}

// 以下是各日志级别的便捷调用函数
func Info(format string, a ...interface{})     { Log(INFO, format, a...) }
func Debug(format string, a ...interface{})    { Log(DEBUG, format, a...) }
func Warn(format string, a ...interface{})     { Log(WARN, format, a...) }
func Error(format string, a ...interface{})    { Log(ERROR, format, a...) }
func VVerbose(format string, a ...interface{}) { Log(VVERBOSE, format, a...) }

// GetEnv 获取环境变量，如果不存在则返回备选默认值
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// Mask 用于脱敏处理敏感字符串（如 API 密钥、密码等）
func Mask(s string) string {
	runes := []rune(s)
	n := len(runes)

	// 长度太短直接全部屏蔽
	if n <= 4 {
		return "****"
	}

	// 计算展示的字符数量，最多显示 10 个字符或长度的三分之一
	showCount := 10
	if n <= showCount {
		showCount = n / 3
	}

	return string(runes[:showCount]) + "[已脱敏]"
}

// CheckEnvFile 检查指定的 .env 配置文件是否存在
func CheckEnvFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Error("严重错误：未在路径 %s 找到 .env 配置文件", filePath)
		return
	}

	Warn("确认：.env 配置文件存在于路径 %s", filePath)
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
