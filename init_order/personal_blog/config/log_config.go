package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"personal_blog/global"
)

// InitLogger 初始化zap日志库
func InitLogger() {
	// 创建日志目录
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,  // 每个日志文件最大100MB
		MaxBackups: 3,    // 最多保留3个备份文件
		MaxAge:     7,    // 最多保留7天
		Compress:   true, // 压缩旧文件
	}

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建写入器
	fileWriter := zapcore.AddSync(logWriter)
	consoleWriter := zapcore.Lock(os.Stdout)

	// 创建核心
	fileCore := zapcore.NewCore(encoder, fileWriter, zapcore.InfoLevel)
	consoleCore := zapcore.NewCore(encoder, consoleWriter, zapcore.InfoLevel)

	// 创建组合核心
	core := zapcore.NewTee(fileCore, consoleCore)

	// 创建logger
	logger := zap.New(core, zap.AddCaller())

	// 将logger设置到全局变量中
	global.Logger = logger

	// 记录初始化日志
	global.Logger.Info("Logger initialized successfully with log rotation")
}
