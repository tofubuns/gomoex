// package xzap 通过对 zap 的二次封装来实现了, 更常用简单的使用, 一致化的日志风格.
package xzap

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger 原始包 zap.Logger 的类型别名
	Logger = zap.Logger
	// Option 原始包 zap.Option 的类型别名
	Option = zap.Option
	// Level 原始包 zapcore.Level 的类型别名
	Level = zapcore.Level
)

// init 初始化开箱即用的日志实例
func init() {
	New(os.Stdout, DebugLevel, AddCaller(), AddStacktrace(DebugLevel))
}

// New 根据指定的参数, 创建一个日志实例, 并将其替换为全局日志实例
//
// 创建模式为单例模式, 不可以多次使用此函数创建日志实例, 想要获取新的实例应该使用
// 日志对象上的 With/WithOptions 方法来获取, 否则 新的日志实例将会覆盖旧的实例
func New(out io.Writer, level Level, opts ...Option) *Logger {
	logger := zap.New(NewCore(out, zap.NewAtomicLevelAt(level)), opts...)
	zap.ReplaceGlobals(logger)
	return zap.L()
}

// Tree 定义了一个输出和启用器函数
type Tree struct {
	Out     io.Writer        // 输出句柄
	Enabler func(Level) bool // 启用器
}

// NewTree 根据指定的参数, 创建一个多核心的日志实例, 每一个Tree结果都对应一个日志核心
//
// 创建模式为单例模式, 不可以多次使用此函数创建日志实例, 想要获取新的实例应该使用
// 日志对象上的 With/WithOptions 方法来获取, 否则 新的日志实例将会覆盖旧的实例
func NewTree(trees []Tree, opts ...Option) *Logger {
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

// Global 得到全局日志实例
func Global() *zap.Logger { return zap.L() }

var (
	// RelpaceGlobals 替换全局日志对象
	RelpaceGlobals = zap.ReplaceGlobals
	// ParseLevel 解析字符串类型的日志等级
	ParseLevel = zapcore.ParseLevel

	// DebugLevel 调试日志级别, 只在调试/开发中使用此级别
	DebugLevel = zapcore.DebugLevel
	// InfoLevel 信息日志等级, 普通等级日志, 通常只记录一些普通信息
	InfoLevel = zapcore.InfoLevel
	// WarnLevel 警告日志等级, 比 Info 等级稍高, 通常需要处理, 但不需要立即处理
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel 错误日志等级, 严重影响用户提要或者系统出现问题使用, 需要立刻处理
	ErrorLevel = zapcore.ErrorLevel
	// PanicLevel 恐慌日志等级, 会导致系统恐慌, 需要立即处理
	PanicLevel = zapcore.PanicLevel
	// FatalLevel 致命日志等级, 系统出现了严重错误, 会停止向外提供服务
	FatalLevel = zapcore.FatalLevel

	// WrapCore 可选配置, 可以包装/替换现有的日志对象 (zap.Logger) 核心 (zapcore.Core)
	WrapCore = zap.WrapCore
	// Hooks 可选配置 回调钩子, 每次输出都会触发注册的钩子
	Hooks = zap.Hooks
	// Fields 可选配置, 增加额外字段, 会在最终的日志中增加额外的字段
	Fields = zap.Fields
	// ErrorOutput 可循啊配置, 错误配置, 会对错误的日志 (Warn 等以上等级) 进行处理
	ErrorOutput = zap.ErrorOutput
	// Development 可选配置, 会将日志设置为开发模式
	Development = zap.Development
	// AddCaller 可选参数 开启调用者信息记录
	AddCaller = zap.AddCaller
	// WithCaller 可选参数, 添加调用者的文件名、行号和函数名称到日志信息中
	WithCaller = zap.WithCaller
	// AddCallerSkip 可选参数, 指定调用者的探查深度
	AddCallerSkip = zap.AddCallerSkip
	// AddStacktrace 可选参数, 指定从什么等级开始记录日志堆栈信息
	AddStacktrace = zap.AddStacktrace
	// IncreaseLevel 可选参数, 提升日志记录器的等级, 只会提升
	IncreaseLevel = zap.IncreaseLevel
	// WithFatalHook 可选参数, 致命错误的钩子
	WithFatalHook = zap.WithFatalHook
	// WithClock 可选参数, 指定日志时钟
	WithClock = zap.WithClock

	Error       = zap.Error
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Object      = zap.Object
	Inline      = zap.Inline
	Any         = zap.Any
)
