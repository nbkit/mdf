package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Flow struct {
	buf   []Field
	log   Logger
	level zapcore.Level
}

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Flow{
			buf: make([]Field, 0, 10),
		}
	},
}

func putEvent(e *Flow) {
	const maxSize = 1 << 16 // 64KiB
	if cap(e.buf) > maxSize {
		return
	}
	eventPool.Put(e)
}
func newFlow(level zapcore.Level, log Logger) *Flow {
	e := eventPool.Get().(*Flow)
	e.buf = e.buf[:0]
	e.log = log
	e.level = level
	return e
}

func (f *Flow) write(msg string) {
	switch f.level {
	case zap.InfoLevel:
		f.log.Info(msg, f.buf...)
		break
	case zap.WarnLevel:
		f.log.Warn(msg, f.buf...)
		break
	case zap.DebugLevel:
		f.log.Debug(msg, f.buf...)
		break
	case zap.FatalLevel:
		f.log.Fatal(msg, f.buf...)
		break
	case zap.ErrorLevel:
		f.log.Error(msg, f.buf...)
		break
	default:
		f.log.Error(msg, f.buf...)
		break
	}
	putEvent(f)
}
func (f *Flow) Msg(msg string) {
	if f == nil {
		return
	}
	f.write(msg)
}

// Send is equivalent to calling Msg("").
//
// NOTICE: once this method is called, the *Event should be disposed.
func (f *Flow) Send() {
	if f == nil {
		return
	}
	f.write("")
}

// Msgf sends the event with formatted msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *Event should be disposed.
// Calling Msgf twice can have unexpected result.
func (f *Flow) Msgf(format string, v ...interface{}) {
	if f == nil {
		return
	}
	f.write(fmt.Sprintf(format, v...))
}
func (f *Flow) Enabled() bool {
	return f != nil && f.log.getLevel().Enabled(f.level)
}

func (f *Flow) Bool(key string, val bool) *Flow {
	f.buf = append(f.buf, zap.Bool(key, val))
	return f
}
func (f *Flow) Float64(key string, val float64) *Flow {
	f.buf = append(f.buf, zap.Float64(key, val))
	return f
}
func (f *Flow) Float32(key string, val float32) *Flow {
	f.buf = append(f.buf, zap.Float32(key, val))
	return f
}

func (f *Flow) Int(key string, val int) *Flow {
	f.buf = append(f.buf, zap.Int(key, val))
	return f
}

func (f *Flow) Int64(key string, val int64) *Flow {
	f.buf = append(f.buf, zap.Int64(key, val))
	return f
}

func (f *Flow) Int32(key string, val int32) *Flow {
	f.buf = append(f.buf, zap.Int32(key, val))
	return f
}

func (f *Flow) Int16(key string, val int16) *Flow {
	f.buf = append(f.buf, zap.Int16(key, val))
	return f
}

func (f *Flow) Int8(key string, val int8) *Flow {
	f.buf = append(f.buf, zap.Int8(key, val))
	return f
}

func (f *Flow) String(key string, val string) *Flow {
	f.buf = append(f.buf, zap.String(key, val))
	return f
}

func (f *Flow) Uint(key string, val uint) *Flow {
	f.buf = append(f.buf, zap.Uint(key, val))
	return f
}

func (f *Flow) Uint64(key string, val uint64) *Flow {
	f.buf = append(f.buf, zap.Uint64(key, val))
	return f
}

func (f *Flow) Uint32(key string, val uint32) *Flow {
	f.buf = append(f.buf, zap.Uint32(key, val))
	return f
}

func (f *Flow) Uint16(key string, val uint16) *Flow {
	f.buf = append(f.buf, zap.Uint16(key, val))
	return f
}

func (f *Flow) Uint8(key string, val uint8) *Flow {
	f.buf = append(f.buf, zap.Uint8(key, val))
	return f
}

func (f *Flow) Time(key string, val time.Time) *Flow {
	f.buf = append(f.buf, zap.Time(key, val))
	return f
}

func (f *Flow) Duration(key string, val time.Duration) *Flow {
	f.buf = append(f.buf, zap.Duration(key, val))
	return f
}
func (f *Flow) Any(key string, val interface{}) *Flow {
	f.buf = append(f.buf, zap.Any(key, val))
	return f
}
