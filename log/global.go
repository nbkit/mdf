package log

import (
	"github.com/nbkit/mdf/internal/zap"
	"github.com/nbkit/mdf/internal/zap/zapcore"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Level = zap.AtomicLevel
type Field = zap.Field

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zap.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

type logConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
	Stack bool   `mapstructure:"stack"`
	debug bool   `mapstructure:"debug"`
	v     *viper.Viper
}

func readConfig() *logConfig {
	config := &logConfig{}
	v := viper.New()
	v.SetConfigName("app")
	v.AddConfigPath(joinCurrentPath("env"))
	if err := v.ReadInConfig(); err != nil {
	}
	if err := v.UnmarshalKey("log", config); err != nil {
	}
	if config.Path == "" {
		config.Path = "./storage/logs"
	}
	config.v = v
	return config
}
func getLevelByTag(tag string) zapcore.Level {
	switch tag {
	case zap.DebugLevel.String():
		return zap.DebugLevel
	case zap.InfoLevel.String():
		return zap.InfoLevel
	case zap.WarnLevel.String():
		return zap.WarnLevel
	case zap.ErrorLevel.String():
		return zap.ErrorLevel
	case zap.FatalLevel.String():
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
	return zap.InfoLevel
}
func pathExists(path string) bool {
	path = joinCurrentPath(path)
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func joinCurrentPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if path.IsAbs(p) {
		return p
	}
	return path.Join(getCurrentPath(), p)
}
func getCurrentPath() string {
	dir := filepath.Dir(os.Args[0])
	dir, _ = filepath.Abs(dir)
	return strings.Replace(dir, "\\", "/", -1)
}
func getFilePath(config *logConfig, args ...string) string {
	parts := make([]string, 0)
	parts = append(parts, joinCurrentPath(config.Path))
	parts = append(parts, "/")
	parts = append(parts, time.Now().Format("20060102"))
	if args != nil && len(args) > 0 {
		parts = append(parts, args...)
	}
	parts = append(parts, ".log")
	return strings.Join(parts, "")
}
