package commonimport

import (
	"admin-panel/pkg/resources"
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"time"
)

var setUpLoggerOnce sync.Once
var logManager Logger

type Logger struct {
	record zerolog.Logger
}

func (l Logger) LogError(ctx context.Context, err error) {
	if err != nil {
		childLogger := l.record.With().Logger()
		if ctx != nil {
			data, ok := ctx.Value(resources.Data).(map[string]interface{})

			if ok {
				for key, value := range data {
					childLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
						return c.Interface(key, value)
					})
				}
			}

		}
		childLogger.Error().Err(err).Send()
	}
}

func (l Logger) LogInfo(ctx context.Context, info string) {
	childLogger := l.record.With().Logger()
	if ctx != nil {
		data, ok := ctx.Value(resources.Data).(map[string]interface{})

		if ok {
			for key, value := range data {
				childLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Interface(key, value)
				})
			}
		}

	}
	childLogger.Info().Str("info", info).Send()
}

func (l Logger) LogDebug(ctx context.Context, debug string) {
	childLogger := l.record.With().Logger()
	if ctx != nil {
		data, ok := ctx.Value(resources.Data).(map[string]interface{})

		if ok {
			for key, value := range data {
				childLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Interface(key, value)
				})
			}
		}

	}
	childLogger.Debug().Str("debug", debug).Send()
}

func (l Logger) LogWarning(ctx context.Context, warning string) {
	childLogger := l.record.With().Logger()
	if ctx != nil {
		data, ok := ctx.Value(resources.Data).(map[string]interface{})

		if ok {
			for key, value := range data {
				childLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Interface(key, value)
				})
			}
		}

	}
	childLogger.Warn().Str("warning", warning).Send()
}

func (l Logger) AddDataToLog(ctx context.Context, data map[string]interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	var logData = make(map[string]interface{})
	buffer, ok := ctx.Value(resources.Data).(map[string]interface{})
	if ok {
		logData = buffer
	}
	for key, value := range data {
		logData[key] = value
	}
	return context.WithValue(ctx, resources.Data, logData)
}

func ProvideLogger() *Logger {
	setUpLoggerOnce.Do(func() {
		fileLogger := &lumberjack.Logger{
			Filename:   "./logs/log.log",
			MaxSize:    50, //
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}
		writers := []io.Writer{fileLogger}
		if true {
			writers = append(writers, os.Stdout)
		}
		output := zerolog.MultiLevelWriter(writers...)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		//logManager.record = zerolog.New(output).With().Timestamp().Stack().Logger()
		logManager.record = zerolog.New(output).With().Timestamp().Logger()
	})
	return &logManager
}
