package log

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

var mapLogs map[string]Logger = make(map[string]Logger)
var mu sync.Mutex

type logger struct {
	level zap.AtomicLevel
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	oracleRegexp             = regexp.MustCompile(`\:\d+`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
	Stack bool   `mapstructure:"stack"`
	debug bool   `mapstructure:"debug"`
}

func readConfig() *LogConfig {
	config := &LogConfig{}
	viper.SetConfigType("yaml")

	viper.SetConfigName("app")
	viper.AddConfigPath(joinCurrentPath("env"))
	if err := viper.ReadInConfig(); err != nil {
		//Errorf("Fatal error when reading %s config file:%s", "app", err)
	}
	if err := viper.UnmarshalKey("log", config); err != nil {
		//Errorf("Fatal error when reading %s config file:%s", "app", err)
	}
	if config.Path == "" {
		config.Path = "./storage/logs"
	}
	return config
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
func getLevelByTag(tag string) zapcore.Level {
	switch tag {
	case zap.DebugLevel.String():
		return zap.DebugLevel
	case zap.InfoLevel.String():
		return zap.InfoLevel
	case zap.ErrorLevel.String():
		return zap.ErrorLevel
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
func getFilePath(config *LogConfig, args ...string) string {
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
func getInstance(key string) Logger {
	mu.Lock()
	defer mu.Unlock()
	if mapLogs[key] == nil {
		mapLogs[key] = createLogger(key)
	}
	return mapLogs[key]
}
func GetLogger(key string) Logger {
	return getInstance(key)
}

func createLogger(args ...string) Logger {
	envConfig := readConfig()
	fileLogger := lumberjack.Logger{
		Filename:   getFilePath(envConfig, args...),
		MaxSize:    10, //MB
		MaxAge:     1,
		MaxBackups: 180,
		LocalTime:  true,
		Compress:   true,
	}
	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(getLevelByTag(envConfig.Level))

	encodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encoderFileConfig := zap.NewProductionEncoderConfig()
	encoderFileConfig.EncodeTime = encodeTime
	encoderFile := zapcore.NewJSONEncoder(encoderFileConfig)

	encoderConsoleConfig := zap.NewDevelopmentEncoderConfig()
	encoderConsoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConsoleConfig.EncodeTime = encodeTime
	encoderConsole := zapcore.NewConsoleEncoder(encoderConsoleConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), logLevel), //打印到控制台
		zapcore.NewCore(encoderFile, zapcore.AddSync(&fileLogger), logLevel),
	)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	l := &logger{
		level: logLevel,
		log:   log,
	}
	l.sugar = l.log.Sugar()
	return l
}

func (s *logger) newEvent(level zapcore.Level, done func(string)) *Flow {
	e := newFlow(level, s)
	return e
}
func (s *logger) Print(v ...interface{}) {
	if v != nil && len(v) > 0 && v[0] == "sql" {
		s.sqlLog(v...)
	} else {
		s.sugar.Debug(v)
	}
}
func (s *logger) CheckAndPrintError(flag string, err error) {
	if err != nil {
		s.Print(flag, err)
	}
}
func (s *logger) SetLevel(tag string) {
	s.level.SetLevel(getLevelByTag(tag))
}

func (s *logger) getLevel() zap.AtomicLevel {
	return s.level
}

//nor

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *logger) Debug(msg string, fields ...Field) {
	s.log.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *logger) Info(msg string, fields ...Field) {
	s.log.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *logger) Warn(msg string, fields ...Field) {
	s.log.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *logger) Error(msg interface{}, fields ...Field) error {
	if ev, ok := msg.(error); ok {
		s.log.Error(ev.Error(), fields...)
		return ev
	} else if ev, ok := msg.(string); ok {
		s.log.Error(ev, fields...)
		return fmt.Errorf(ev)
	}
	s.log.Error(fmt.Sprint(msg), fields...)
	return fmt.Errorf(fmt.Sprint(msg))
}
func (s *logger) Fatal(msg string, fields ...Field) {
	s.log.Fatal(msg, fields...)
}

//f

// Debugf uses fmt.Sprintf to log a templated message.
func (s *logger) DebugF(template string, args ...interface{}) {
	s.sugar.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *logger) InfoF(template string, args ...interface{}) {
	s.sugar.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *logger) WarnF(template string, args ...interface{}) {
	s.sugar.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *logger) ErrorF(template string, args ...interface{}) error {
	s.sugar.Errorf(template, args...)
	return fmt.Errorf(template, args...)
}
func (s *logger) FatalF(template string, args ...interface{}) {
	s.sugar.Fatalf(template, args...)
}

func (s *logger) InfoS() *Flow {
	return s.newEvent(zap.InfoLevel, nil)
}
func (s *logger) ErrorS() *Flow {
	return s.newEvent(zap.ErrorLevel, nil)
}
func (s *logger) WarnS() *Flow {
	return s.newEvent(zap.WarnLevel, nil)
}
func (s *logger) DebugS() *Flow {
	return s.newEvent(zap.DebugLevel, nil)
}

func (s *logger) FatalS() *Flow {
	return s.newEvent(zap.FatalLevel, nil)
}
