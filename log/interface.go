package log

import (
	"go.uber.org/zap"
)

type Logger interface {
	CheckAndPrintError(flag string, err error)
	SetLevel(level string)
	getLevel() zap.AtomicLevel
	Print(v ...interface{})

	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg interface{}, fields ...Field) error
	Fatal(msg string, fields ...Field)

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{}) error
	Fatalf(template string, args ...interface{})

	DebugD() *Flow
	InfoD() *Flow
	WarnD() *Flow
	ErrorD() *Flow
	FatalD() *Flow
}
