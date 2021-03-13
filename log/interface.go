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

	DebugF(template string, args ...interface{})
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{}) error
	FatalF(template string, args ...interface{})

	DebugS() *Flow
	InfoS() *Flow
	WarnS() *Flow
	ErrorS() *Flow
	FatalS() *Flow
}
