package chix

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	l *zap.Logger
}

func fieldsToInterface(fields []Field) []interface{} {
	args := make([]interface{}, len(fields))
	for i, f := range fields {
		args[i] = f
	}
	return args
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.l.Sugar().Debugw(msg, fieldsToInterface(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.l.Sugar().Infow(msg, fieldsToInterface(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.l.Sugar().Warnw(msg, fieldsToInterface(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.l.Sugar().Errorw(msg, fieldsToInterface(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.l.Sugar().Fatalw(msg, fieldsToInterface(fields)...)
}

func NewZapLogger() *ZapLogger {
	logger, _ := zap.NewProduction()
	return &ZapLogger{logger}
}

var defaultLogger Logger = NewZapLogger()
