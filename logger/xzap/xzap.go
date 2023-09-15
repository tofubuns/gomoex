// package xzap 通过对 zap 的二次封装来实现了, 通用的创建, 一致化的日志风格.
package xzap

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// init 初始化开箱即用的日志实例
func init() {
	New(os.Stdout, zap.DebugLevel, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
}

// New 根据指定的参数, 创建一个日志实例, 并将其替换为全局日志实例
//
// 创建模式为单例模式, 不可以多次使用此函数创建日志实例, 想要获取新的实例应该使用
// 日志对象上的 With/WithOptions 方法来获取, 否则 新的日志实例将会覆盖旧的实例
func New(out io.Writer, level zapcore.Level, opts ...zap.Option) *zap.Logger {
	logger := zap.New(NewCore(out, zap.NewAtomicLevelAt(level)), opts...)
	zap.ReplaceGlobals(logger)
	return zap.L()
}

// Tree 定义了一个输出和启用器函数
type Tree struct {
	Out     io.Writer                // 输出句柄
	Enabler func(zapcore.Level) bool // 启用器
}

// NewTree 根据指定的参数, 创建一个多核心的日志实例, 每一个Tree结果都对应一个日志核心
//
// 创建模式为单例模式, 不可以多次使用此函数创建日志实例, 想要获取新的实例应该使用
// 日志对象上的 With/WithOptions 方法来获取, 否则 新的日志实例将会覆盖旧的实例
func NewTree(trees []Tree, opts ...zap.Option) *zap.Logger {
	var cores []zapcore.Core

	for _, tree := range trees {
		cores = append(cores, NewCore(
			tree.Out,
			zap.LevelEnablerFunc(tree.Enabler),
		))
	}

	logger := zap.New(zapcore.NewTee(cores...), opts...)
	zap.ReplaceGlobals(logger)
	return zap.L()
}

// NewCore 创建一个基于 zap 的标准日志核心
func NewCore(out io.Writer, le zapcore.LevelEnabler) zapcore.Core {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.NameKey = "logger"
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.MessageKey = "message"
	encoderCfg.FunctionKey = "function"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.DateTime))
	}

	encoder := zapcore.NewJSONEncoder(encoderCfg)
	syncOut := zapcore.AddSync(out)
	return zapcore.NewCore(encoder, syncOut, le)
}
