package main

import (
	"os"

	"github.com/tofubuns/gomoex/logger/xzap"
)

func main() {
	// 开箱即用
	{
		xzap.Global().Info("test message", xzap.Int("Code", 1))
		xzap.Global().Sync()
	}

	// 配置
	{
		loggerOptions := []xzap.Option{
			xzap.AddCaller(),
			xzap.AddStacktrace(xzap.DebugLevel),
		}

		logger := xzap.New(os.Stdout, xzap.DebugLevel, loggerOptions...)
		logger.Debug("test message", xzap.Int("Code", 2))
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
			xzap.Global().Error(err.Error())
		}
		defer urgentFile.Close()

		normalFile, err := os.OpenFile(normalLogsFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			xzap.Global().Error(err.Error())
		}
		defer normalFile.Close()

		logger := xzap.NewTree([]xzap.Tree{
			{
				Out: urgentFile,
				Enabler: func(level xzap.Level) bool {
					return level >= xzap.ErrorLevel
				},
			},
			{
				Out: normalFile,
				Enabler: func(level xzap.Level) bool {
					return level <= xzap.WarnLevel
				},
			},
		})
		logger.Debug("test message", xzap.Int("Code", 3))
		logger.Error("test message", xzap.Int("Code", 3))
		logger.Sync()
	}
}
