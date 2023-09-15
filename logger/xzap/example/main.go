package main

import (
	"os"

	"github.com/tofubuns/gomoex/logger/xzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 开箱即用
	{
		zap.L().Info("test message", zap.Int("Code", 1))
		zap.L().Sync()
	}

	// 配置
	{
		loggerOptions := []zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zap.DebugLevel),
		}

		logger := xzap.New(os.Stdout, zap.DebugLevel, loggerOptions...)
		logger.Debug("test message", zap.Int("Code", 2))
		logger.Sync()
	}

	// 多位置日志
	{
		var (
			urgentLogsFilename = "urgent.log" // 紧急日志
			normalLogsFilename = "normal.log" // 平常日志
		)

		urgentFile, err := os.OpenFile(urgentLogsFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			zap.L().Error(err.Error())
		}
		defer urgentFile.Close()

		normalFile, err := os.OpenFile(normalLogsFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			zap.L().Error(err.Error())
		}
		defer normalFile.Close()

		logger := xzap.NewTree([]xzap.Tree{
			{
				Out: urgentFile,
				Enabler: func(level zapcore.Level) bool {
					return level >= zap.ErrorLevel
				},
			},
			{
				Out: normalFile,
				Enabler: func(level zapcore.Level) bool {
					return level <= zap.WarnLevel
				},
			},
		})
		logger.Debug("test message", zap.Int("Code", 3))
		logger.Error("test message", zap.Int("Code", 3))
		logger.Sync()
	}
}
