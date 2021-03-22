package log

import (
	"fmt"
	"github.com/nbkit/mdf/internal/zap"
	"github.com/nbkit/mdf/internal/zap/zapcore"
	"sync"
	"time"
)

type Flow struct {
	buf        []Field
	output     IOutput
	level      zapcore.Level
	time       time.Time
	callerSkip int
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
func newFlow(level zapcore.Level, output IOutput) *Flow {
	e := eventPool.Get().(*Flow)
	e.buf = e.buf[:0]
	e.output = output
	e.level = level
	e.time = time.Now()
	e.callerSkip = 0
	return e
}

func (f *Flow) write(msg string, fields ...Field) {
	if len(fields) > 0 {
		f.buf = append(f.buf, fields...)
	}
	switch f.level {
	case zap.InfoLevel:
		f.output.Clone(callerSkip(f.callerSkip)).Info(msg, f.buf...)
		break
	case zap.WarnLevel:
		f.output.Clone(callerSkip(f.callerSkip)).Warn(msg, f.buf...)
		break
	case zap.DebugLevel:
		f.output.Clone(callerSkip(f.callerSkip)).Debug(msg, f.buf...)
		break
	case zap.FatalLevel:
		f.output.Clone(callerSkip(f.callerSkip)).Fatal(msg, f.buf...)
		break
	case zap.ErrorLevel:
		f.output.Clone(callerSkip(f.callerSkip)).Error(msg, f.buf...)
		break
	default:
		f.output.Clone(callerSkip(f.callerSkip)).Error(msg, f.buf...)
		break
	}
	putEvent(f)
}

func (f *Flow) Enabled() bool {
	return f != nil && f.output.GetLevel().Enabled(f.level)
}
func (f *Flow) Output(msg ...interface{}) {
	if f == nil {
		return
	}
	if len(msg) > 0 {
		f.write(fmt.Sprint(msg...))
	} else {
		f.write("")
	}
}
func (f *Flow) Msg(msg string, fields ...Field) {
	if f == nil {
		return
	}
	f.write(msg, fields...)
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

func (f *Flow) CallerSkip(skip int) *Flow {
	f.callerSkip = skip
	return f
}

func (f *Flow) Latency() *Flow {
	f.buf = append(f.buf, zap.Duration(LatencyFieldName, time.Now().Sub(f.time)))
	f.time = time.Now()
	return f
}
func (f *Flow) Error(msg interface{}, fields ...Field) error {
	if ev, ok := msg.(error); ok {
		f.write(ev.Error(), fields...)
		return ev
	} else if ev, ok := msg.(string); ok {
		f.write(ev, fields...)
		return fmt.Errorf(ev)
	}
	f.write(fmt.Sprint(msg), fields...)
	return fmt.Errorf(fmt.Sprint(msg))
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
func (f *Flow) StackSkip(key string, skip int) *Flow {
	f.buf = append(f.buf, zap.StackSkip(key, skip))
	return f
}
