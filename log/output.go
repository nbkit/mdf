package log

import (
	"fmt"
	"github.com/nbkit/mdf/internal/zap"
	"github.com/nbkit/mdf/internal/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

type IOutput interface {
	SetLevel(tag string)
	GetLevel() Level
	Clone(opts ...zap.Option) IOutput
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg interface{}, fields ...Field) error
	Fatal(msg string, fields ...Field)
}

var outputMap map[string]IOutput = make(map[string]IOutput)
var mu sync.Mutex

type DefaultOutput struct {
	level      zap.AtomicLevel
	log        *zap.Logger
	callerSkip int
}

func getOutputInstance(key string) IOutput {
	mu.Lock()
	defer mu.Unlock()
	if outputMap[key] == nil {
		outputMap[key] = createDefaultOutput(key)
	}
	return outputMap[key]
}
func GetOutput(key string) IOutput {
	return getOutputInstance(key)
}
func SetOutput(key string, log IOutput) {
	mu.Lock()
	defer mu.Unlock()
	outputMap[key] = log
}
func createDefaultOutput(args ...string) IOutput {
	envConfig := readConfig()
	fileLogger := lumberjack.Logger{
		Filename:   getFilePath(envConfig, args...),
		MaxSize:    10, //MB
		MaxAge:     1,
		MaxBackups: 180,
		LocalTime:  true,
		Compress:   true,
	}
	levelPath := "log.level"
	logLevel := zap.NewAtomicLevel()
	if len(args) > 0 && args[0] != "" {
		levelPath = levelPath + "." + args[0]
	}
	if l := envConfig.v.GetString(levelPath); l != "" {
		logLevel.SetLevel(getLevelByTag(l))
	} else {
		logLevel.SetLevel(getLevelByTag(envConfig.Level))
	}
	encodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encoderFileConfig := zap.NewProductionEncoderConfig()
	encoderFileConfig.EncodeTime = encodeTime
	encoderFile := zapcore.NewJSONEncoder(encoderFileConfig)

	encoderConsoleConfig := zap.NewDevelopmentEncoderConfig()
	encoderConsoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConsoleConfig.EncodeTime = encodeTime
	encoderConsole := zapcore.NewConsoleEncoder(encoderConsoleConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), logLevel), //打印到控制台
		zapcore.NewCore(encoderFile, zapcore.AddSync(&fileLogger), logLevel),
	)
	opts := make([]zap.Option, 0)
	if envConfig.Stack {
		zap.AddStacktrace(logLevel)
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(3))
	}
	log := zap.New(core, opts...)
	l := &DefaultOutput{
		level: logLevel,
		log:   log,
	}
	return l
}

func (s *DefaultOutput) Clone(opts ...zap.Option) IOutput {
	copy := *s
	copy.log = copy.log.WithOptions(opts...)
	return &copy
}
func (s *DefaultOutput) SetLevel(tag string) {
	s.level.SetLevel(getLevelByTag(tag))
}

func (s *DefaultOutput) GetLevel() Level {
	return s.level
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the output site, as well as any fields accumulated on the DefaultOutput.
func (s *DefaultOutput) Debug(msg string, fields ...Field) {
	s.log.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the output site, as well as any fields accumulated on the DefaultOutput.
func (s *DefaultOutput) Info(msg string, fields ...Field) {
	s.log.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the output site, as well as any fields accumulated on the DefaultOutput.
func (s *DefaultOutput) Warn(msg string, fields ...Field) {
	s.log.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the output site, as well as any fields accumulated on the DefaultOutput.
func (s *DefaultOutput) Error(msg interface{}, fields ...Field) error {
	if ev, ok := msg.(error); ok {
		s.log.Error(ev.Error(), fields...)
		return ev
	} else if ev, ok := msg.(string); ok {
		s.log.Error(ev, fields...)
		return fmt.Errorf(ev)
	}
	s.log.Error(fmt.Sprint(msg), fields...)
	return fmt.Errorf(fmt.Sprint(msg))
}
func (s *DefaultOutput) Fatal(msg string, fields ...Field) {
	s.log.Fatal(msg, fields...)
}
