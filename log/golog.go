package log

// Default is the package-level ready-to-use logger,
// level had set to "info", is changeable.
var Default = createLogger()

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
func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) error {
	return Default.Errorf(format, args...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

func InfoD() *Flow {
	return Default.InfoD()
}
func ErrorD() *Flow {
	return Default.ErrorD()
}
func DebugD() *Flow {
	return Default.DebugD()
}
func FatalD() *Flow {
	return Default.FatalD()
}
