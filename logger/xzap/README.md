## xzap

基于 [zap](https://pkg.go.dev/go.uber.org/zap) 二次封装, 此封装的封装理念为

1. 提供更简单更常用的创建方式
2. 固定日志格式/风格
3. 提供日常常用的 API
4. 以类型别名的形式存在让使用者可以不直接依赖于原始包
5. 提供开箱即用的日志实例
6. 使用全局日志方式实现单例模式
7. 不再提供全局方法日志

### 使用示例

开箱即用
```golang
package main

import (
    "github.com/tofubuns/gocomx/logger/xzap"
)

func main() {
    xzap.Global().Info("test message", xzap.Int("Code", 1))
    xzap.Global().Sync()
}
```

### 原生配置
支持原始包的原生配置
```golang
    loggerOptions := []xzap.Option{
    	xzap.AddCaller(),
    	xzap.AddStacktrace(xzap.DebugLevel),
    }

    logger := xzap.New(os.Stdout, xzap.DebugLevel, loggerOptions...)
    logger.Debug("test message", xzap.Int("Code", 2))
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
```
