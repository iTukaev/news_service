package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

var (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
)
var levelMap = map[LogLevel]zapcore.Level{
	Debug: zapcore.DebugLevel,
	Info:  zapcore.InfoLevel,
	Error: zapcore.ErrorLevel,
	Fatal: zapcore.FatalLevel,
}

func getLoggerLevel(lvl LogLevel) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func New(lvl LogLevel) (*zap.SugaredLogger, error) {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(getLoggerLevel(lvl)),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "lvl",
			TimeKey:        "time",
			FunctionKey:    "file",
			NameKey:        "log",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build new logger")
	}

	return logger.Sugar(), nil
}
