package log

import (
	"fmt"
	"sync"
)

// defaultOutput is the package-level ready-to-use DefaultOutput,
// level had set to "info", is changeable.
var (
	LatencyFieldName = "latency"
	defaultLog       = Get("")
)

type ILog interface {
	Warn() *Flow
	Info() *Flow
	Error() *Flow
	Debug() *Flow
	Fatal() *Flow
	FatalF(format string, args ...interface{})
	ErrorF(format string, args ...interface{}) error
	WarnF(format string, args ...interface{})
	InfoF(format string, args ...interface{})
	DebugF(format string, args ...interface{})
}

var logMap map[string]ILog = make(map[string]ILog)
var logmu sync.Mutex

type logImpl struct {
	output IOutput
}

func Get(key string) ILog {
	logmu.Lock()
	defer logmu.Unlock()
	if logMap[key] == nil {
		logMap[key] = createLogger(key)
	}
	return logMap[key]
}
func createLogger(key string) ILog {
	item := &logImpl{
		output: getOutputInstance(key),
	}
	return item
}
func (f *logImpl) Warn() *Flow {
	return newFlow(WarnLevel, f.output)
}
func (f *logImpl) Info() *Flow {
	return newFlow(InfoLevel, f.output)
}
func (f *logImpl) Error() *Flow {
	return newFlow(ErrorLevel, f.output)
}
func (f *logImpl) Debug() *Flow {
	return newFlow(DebugLevel, f.output)
}
func (f *logImpl) Fatal() *Flow {
	return newFlow(FatalLevel, f.output)
}
func (f *logImpl) FatalF(format string, args ...interface{}) {
	f.Fatal().CallerSkip(1).Msgf(format, args...)
}
func (f *logImpl) ErrorF(format string, args ...interface{}) error {
	return f.Error().CallerSkip(1).Error(fmt.Sprintf(format, args...))
}
func (f *logImpl) WarnF(format string, args ...interface{}) {
	f.Warn().CallerSkip(1).Msgf(format, args...)
}
func (f *logImpl) InfoF(format string, args ...interface{}) {
	f.Info().CallerSkip(1).Msgf(format, args...)
}
func (f *logImpl) DebugF(format string, args ...interface{}) {
	f.Debug().CallerSkip(1).Msgf(format, args...)
}

// Fatal `os.Exit(1)` exit no matter the level of the DefaultOutput.
// If the DefaultOutput's level is fatal, error, warn, info or debug
// then it will print the output message too.
func FatalD(msg string, fields ...Field) {
	defaultLog.Fatal().CallerSkip(1).Msg(msg, fields...)
}

// Error will print only when DefaultOutput's Level is error, warn, info or debug.
func ErrorD(msg interface{}, fields ...Field) error {
	return defaultLog.Error().CallerSkip(1).Error(msg, fields...)

}

// Warn will print when DefaultOutput's Level is warn, info or debug.
func WarnD(msg string, fields ...Field) {
	defaultLog.Warn().CallerSkip(1).Msg(msg, fields...)
}

// Info will print when DefaultOutput's Level is info or debug.
func InfoD(msg string, fields ...Field) {
	defaultLog.Info().CallerSkip(1).Msg(msg, fields...)
}

// Debug will print when DefaultOutput's Level is debug.
func DebugD(msg string, fields ...Field) {
	defaultLog.Debug().CallerSkip(1).Msg(msg, fields...)
}

// Fatalf will `os.Exit(1)` no matter the level of the DefaultOutput.
// If the DefaultOutput's level is fatal, error, warn, info or debug
// then it will print the output message too.
func FatalF(format string, args ...interface{}) {
	defaultLog.Fatal().CallerSkip(1).Msgf(format, args...)
}

// Errorf will print only when DefaultOutput's Level is error, warn, info or debug.
func ErrorF(format string, args ...interface{}) error {
	return defaultLog.Error().CallerSkip(1).Error(fmt.Sprintf(format, args...))
}

// Warnf will print when DefaultOutput's Level is warn, info or debug.
func WarnF(format string, args ...interface{}) {
	defaultLog.Warn().CallerSkip(1).Msgf(format, args...)
}

// Infof will print when DefaultOutput's Level is info or debug.
func InfoF(format string, args ...interface{}) {
	defaultLog.Info().CallerSkip(1).Msgf(format, args...)
}

// Debugf will print when DefaultOutput's Level is debug.
func DebugF(format string, args ...interface{}) {
	defaultLog.Debug().CallerSkip(1).Msgf(format, args...)
}
func Print(v ...interface{}) {
	defaultLog.Info().Output(v...)
}
func Warn() *Flow {
	return defaultLog.Warn()
}
func Info() *Flow {
	return defaultLog.Info()
}
func Error() *Flow {
	return defaultLog.Error()
}
func Debug() *Flow {
	return defaultLog.Debug()
}
func Fatal() *Flow {
	return defaultLog.Fatal()
}
