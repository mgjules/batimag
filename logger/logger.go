package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	initialSamplingSec    = 100
	thereafterSamplingSec = 100
)

// Logger is a simple wrapper around zap.SugaredLogger.
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new Logger.
func New(debug bool) (*Logger, error) {
	cfg := zap.Config{
		DisableStacktrace: true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
	}

	if debug {
		cfg.Development = true
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Encoding = "console"
		cfg.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		cfg.Development = false
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

		cfg.Sampling = &zap.SamplingConfig{
			Initial:    initialSamplingSec,
			Thereafter: thereafterSamplingSec,
		}
		cfg.Encoding = "json"
		cfg.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "line",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	var (
		logger *zap.Logger
		err    error
	)
	logger, err = cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("new zap logger: %w", err)
	}

	return &Logger{logger.Sugar()}, nil
}
