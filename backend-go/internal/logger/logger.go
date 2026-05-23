package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志配置
type Config struct {
	// 日志目录
	LogDir string
	// 日志文件名
	LogFile string
	// 单个日志文件最大大小 (MB)
	MaxSize int
	// 保留的旧日志文件最大数量
	MaxBackups int
	// 保留的旧日志文件最大天数
	MaxAge int
	// 是否压缩旧日志文件
	Compress bool
	// 是否同时输出到控制台
	Console bool
}

// rawFileLog 仅写文件的 logger，用于向日志文件写入原始 JSON 输出
var rawFileLog *log.Logger

// consoleLog 仅写控制台的 logger，用于精简格式输出（不写入文件，避免与 raw 日志重复）
var consoleLog *log.Logger

// RawFileLog 返回仅写文件的 logger。
// 未初始化时回退到全局 logger。
func RawFileLog() *log.Logger {
	if rawFileLog != nil {
		return rawFileLog
	}
	return log.Default()
}

// ConsoleLog 返回仅写控制台的 logger。
// 未初始化时回退到全局 logger。
func ConsoleLog() *log.Logger {
	if consoleLog != nil {
		return consoleLog
	}
	return log.Default()
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		LogDir:     "logs",
		LogFile:    "app.log",
		MaxSize:    100, // 100MB
		MaxBackups: 10,
		MaxAge:     30, // 30 days
		Compress:   true,
		Console:    true,
	}
}

// Setup 初始化日志系统
func Setup(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 确保日志目录存在
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	logPath := filepath.Join(cfg.LogDir, cfg.LogFile)

	// 配置 lumberjack 日志轮转
	lumberLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	flags := log.Ldate | log.Ltime | log.Lmicroseconds

	// log.Printf 写入 stdout + 文件（普通日志，如初始化、调度信息等）
	if cfg.Console {
		log.SetOutput(io.MultiWriter(os.Stdout, lumberLogger))
	} else {
		log.SetOutput(lumberLogger)
	}
	// rawFileLog 始终仅写文件（原始 JSON），用于双通道输出
	rawFileLog = log.New(lumberLogger, "", flags)
	// consoleLog 始终仅写控制台（精简格式），用于请求/响应日志的控制台通道
	// 避免精简格式重复写入文件
	if cfg.Console {
		consoleLog = log.New(os.Stdout, "", flags)
	} else {
		// 无控制台模式下，consoleLog 回退到 rawFileLog（仅文件）
		consoleLog = rawFileLog
	}

	log.SetFlags(flags)

	log.Printf("[Logger-Init] 日志系统已初始化")
	log.Printf("[Logger-Init] 日志文件: %s", logPath)
	log.Printf("[Logger-Init] 轮转配置: 最大 %dMB, 保留 %d 个备份, %d 天", cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	return nil
}
