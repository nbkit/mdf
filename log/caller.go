package log

import "github.com/nbkit/mdf/internal/zap"

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*zap.Logger)

func (f optionFunc) Apply(log *zap.Logger) {
	f(log)
}

func CallerSkip(skip int) zap.Option {
	return optionFunc(func(log *zap.Logger) {
		log.CallerSkip += skip
	})
}
