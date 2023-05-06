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

var outputMap map[string]*zap.Logger = make(map[string]*zap.Logger)
var mu sync.Mutex

func (f *Flow) getWriter() *zap.Logger {
	mu.Lock()
	defer mu.Unlock()
	key := fmt.Sprintf("%v:%v:%v:%v", f.name, f.forceOutput, f.level.String(), f.callerSkip)
	if outputMap[key] == nil {
		outputMap[key] = f.createWriter(f.name)
	}
	return outputMap[key]
}
func (f *Flow) createWriter(args ...string) *zap.Logger {
	envConfig := readConfig()
	fileLogger := lumberjack.Logger{
		Filename:   getFilePath(envConfig, args...),
		MaxSize:    10, //MB
		MaxAge:     1,
		MaxBackups: 180,
		LocalTime:  true,
		Compress:   true,
	}
	logLevel := zap.NewAtomicLevel()

	levelPath := "log.level"
	if len(args) > 0 && args[0] != "" {
		levelPath = "log.level." + args[0]
	}
	//log.level.xxx -> log.level.root -> log.level
	if l := envConfig.v.GetString(levelPath); l != "" {
		logLevel.SetLevel(getLevelByTag(l))
	} else if l := envConfig.v.GetString("log.level.root"); l != "" {
		logLevel.SetLevel(getLevelByTag(l))
	} else if l := envConfig.v.GetString("log.level"); l != "" {
		logLevel.SetLevel(getLevelByTag(l))
	} else {
		logLevel.SetLevel(zapcore.ErrorLevel)
	}
	if f.forceOutput {
		logLevel.SetLevel(zapcore.InfoLevel)
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
	opts = append(opts, zap.AddCallerSkip(2))
	if envConfig.Stack {
		zap.AddStacktrace(logLevel)
		opts = append(opts, zap.AddCaller())
	}
	if f.callerSkip > 0 {
		opts = append(opts, zap.AddCallerSkip(f.callerSkip))
	}

	log := zap.New(core, opts...)
	return log
}
