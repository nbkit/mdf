package log

import "fmt"

// defaultOutput is the package-level ready-to-use DefaultOutput,
// level had set to "info", is changeable.
var (
	LatencyFieldName = "latency"
	defaultOutput    = createDefaultOutput()
)

// Fatal `os.Exit(1)` exit no matter the level of the DefaultOutput.
// If the DefaultOutput's level is fatal, error, warn, info or debug
// then it will print the output message too.
func FatalD(msg string, fields ...Field) {
	Fatal().CallerSkip(1).Msg(msg, fields...)
}

// Error will print only when DefaultOutput's Level is error, warn, info or debug.
func ErrorD(msg interface{}, fields ...Field) error {
	return Error().CallerSkip(1).Error(msg, fields...)

}

// Warn will print when DefaultOutput's Level is warn, info or debug.
func WarnD(msg string, fields ...Field) {
	Warn().CallerSkip(1).Msg(msg, fields...)
}

// Info will print when DefaultOutput's Level is info or debug.
func InfoD(msg string, fields ...Field) {
	Info().CallerSkip(1).Msg(msg, fields...)
}

// Debug will print when DefaultOutput's Level is debug.
func DebugD(msg string, fields ...Field) {
	Debug().CallerSkip(1).Msg(msg, fields...)
}

// Fatalf will `os.Exit(1)` no matter the level of the DefaultOutput.
// If the DefaultOutput's level is fatal, error, warn, info or debug
// then it will print the output message too.
func FatalF(format string, args ...interface{}) {
	Fatal().CallerSkip(1).Msgf(format, args...)
}

// Errorf will print only when DefaultOutput's Level is error, warn, info or debug.
func ErrorF(format string, args ...interface{}) error {
	return Error().CallerSkip(1).Error(fmt.Sprintf(format, args...))
}

// Warnf will print when DefaultOutput's Level is warn, info or debug.
func WarnF(format string, args ...interface{}) {
	Warn().CallerSkip(1).Msgf(format, args...)
}

// Infof will print when DefaultOutput's Level is info or debug.
func InfoF(format string, args ...interface{}) {
	Info().CallerSkip(1).Msgf(format, args...)
}

// Debugf will print when DefaultOutput's Level is debug.
func DebugF(format string, args ...interface{}) {
	Debug().CallerSkip(1).Msgf(format, args...)
}
func Print(v ...interface{}) {
	Info().CallerSkip(1).Output(v...)
}
func Warn() *Flow {
	return newFlow(WarnLevel, defaultOutput)
}
func Info() *Flow {
	return newFlow(InfoLevel, defaultOutput)
}
func Error() *Flow {
	return newFlow(ErrorLevel, defaultOutput)
}
func Debug() *Flow {
	return newFlow(DebugLevel, defaultOutput)
}
func Fatal() *Flow {
	return newFlow(FatalLevel, defaultOutput)
}
