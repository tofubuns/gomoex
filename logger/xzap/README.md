## xzap

基于 [zap](https://pkg.go.dev/go.uber.org/zap) 二次封装, 此封装的封装理念为

1. 提供更常用的创建方式
2. 固定日志格式/风格
5. 初始化开箱即用的日志实例

### 使用示例

开箱即用
```golang
package main

import (
    "github.com/tofubuns/gocomx/logger/xzap"
)

func main() {
    zap.L().Info("test message", zap.Int("Code", 1))
    zap.L().Sync()
}
```

### 原生配置
支持原始包的原生配置
```golang
    loggerOptions := []zap.Option{
    	zap.AddCaller(),
    	zap.AddStacktrace(zapcore.DebugLevel),
    }

    logger := xzap.New(os.Stdout, zapcore.DebugLevel, loggerOptions...)
    logger.Debug("test message", zap.Int("Code", 2))
    logger.Sync()
```

### 多日志位置输出
让不同的日志输出道不同的位置/路径
```golang
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
    			return level >= zapcore.ErrorLevel
    		},
    	},
    	{
    		Out: normalFile,
    		Enabler: func(level zapcore.Level) bool {
    			return level <= zapcore.WarnLevel
    		},
    	},
    })
    logger.Debug("test message", zap.Int("Code", 3))
    logger.Error("test message", zap.Int("Code", 3))
    logger.Sync()
```
