package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerOnce sync.Once
var logger *zap.Logger

type Config struct {
	Filename string `json:"filename"`
}

func defaultConfig() *Config {
	return &Config{
		"log/noti.log",
	}
}

func Init(cfg *Config) {
	loggerOnce.Do(func() {
		logger = newLogger(cfg)
	})
}

func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		logger = newLogger(defaultConfig())
	})

	return logger
}

func newLogger(cfg *Config) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.MillisDurationEncoder
	encoderCfg.ConsoleSeparator = "|"
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := newCore(encoder, cfg)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel))
	return log
}

func newCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	syncer := &lumberjack.Logger{
		Filename: cfg.Filename,
	}
	f := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l > zapcore.DebugLevel
	})

	writer := zapcore.AddSync(syncer)
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		f,
	)
	return core
}
