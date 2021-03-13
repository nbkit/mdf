package log

// Default is the package-level ready-to-use logger,
// level had set to "info", is changeable.
var (
	LatencyFieldName = "latency"
	Default          = createLogger()
)

// Reset re-sets the default logger to an empty one.
func Reset() {
	Default = createLogger()
}

// Print prints a log message without levels and colors.
func Print(v ...interface{}) {
	Default.Print(v...)
}
func CheckAndPrintError(flag string, err error) {
	Default.CheckAndPrintError(flag, err)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(msg string, fields ...Field) {
	Default.Fatal(msg, fields...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(msg interface{}, fields ...Field) error {
	return Default.Error(msg, fields...)

}

// Warn will print when logger's Level is warn, info or debug.
func Warn(msg string, fields ...Field) {
	Default.Warn(msg, fields...)
}

// Info will print when logger's Level is info or debug.
func Info(msg string, fields ...Field) {
	Default.Info(msg, fields...)
}

// Debug will print when logger's Level is debug.
func Debug(msg string, fields ...Field) {
	Default.Debug(msg, fields...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func FatalF(format string, args ...interface{}) {
	Default.FatalF(format, args...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func ErrorF(format string, args ...interface{}) error {
	return Default.ErrorF(format, args...)
}

// Warnf will print when logger's Level is warn, info or debug.
func WarnF(format string, args ...interface{}) {
	Default.WarnF(format, args...)
}

// Infof will print when logger's Level is info or debug.
func InfoF(format string, args ...interface{}) {
	Default.InfoF(format, args...)
}

// Debugf will print when logger's Level is debug.
func DebugF(format string, args ...interface{}) {
	Default.DebugF(format, args...)
}

func InfoS() *Flow {
	return Default.InfoS()
}
func ErrorS() *Flow {
	return Default.ErrorS()
}
func DebugS() *Flow {
	return Default.DebugS()
}
func FatalS() *Flow {
	return Default.FatalS()
}
