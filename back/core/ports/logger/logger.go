package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Field struct {
	Key   string
	Value string
}

type Logger interface {
	Info(message string, fields ...Field)
	Debug(message string, fields ...Field)
	Error(err error, message string, fields ...Field)
	WithContext(ctx context.Context) context.Context
}

type logger struct {
	internalLogger *zerolog.Logger
}

func NewLogger(logLevel string) (Logger, error) {
	zLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	interanlLogger := zerolog.New(os.Stderr).Level(zLevel).With().Timestamp().Logger()
	return &logger{
		internalLogger: &interanlLogger,
	}, nil
}

func (l *logger) Info(message string, fields ...Field) {
	addFields(l.internalLogger.Info()).Msg(message)
}

func (l *logger) Debug(message string, fields ...Field) {
	addFields(l.internalLogger.Debug()).Msg(message)
}

func (l *logger) Error(err error, message string, fields ...Field) {
	addFields(l.internalLogger.Error()).Err(err).Msg(message)
}

func addFields(logEvent *zerolog.Event, fields ...Field) *zerolog.Event {
	for _, field := range fields {
		logEvent.Str(field.Key, field.Value)
	}
	return logEvent
}

func Error(ctx context.Context, err error, message string, fields ...Field) {
	ctxLogger(ctx).Error(err, message, fields...)
}

func (l logger) WithContext(ctx context.Context) context.Context {
	return l.internalLogger.WithContext(ctx)
}

func ctxLogger(ctx context.Context) Logger {
	return &logger{internalLogger: zerolog.Ctx(ctx)}
}
