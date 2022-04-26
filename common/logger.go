package common

import (
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerReservedKeyExpected = "Expected"
	LoggerReservedKeyActual   = "Actual"
)

type LoggerFactory interface {
	CreateCommonLogger() Logger
	CreateErrorLogger() Logger
	CreateSuccessLogger() Logger
}

func NewFxLoggerFactory() LoggerFactory {
	return &loggerFactory{}
}

type loggerFactory struct {
}

func (f *loggerFactory) CreateCommonLogger() Logger {
	return f.createLogger(&logger{
		groupHeaderColor:    color.New(color.FgBlue, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgMagenta, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgWhite),
	})
}

func (f *loggerFactory) CreateErrorLogger() Logger {
	return f.createLogger(&logger{
		groupHeaderColor:    color.New(color.FgRed, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgRed, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgRed),
		paramColor2:         color.New(color.FgGreen),
		paramColor3:         color.New(color.FgYellow),
	})
}

func (f *loggerFactory) CreateSuccessLogger() Logger {
	return f.createLogger(&logger{
		groupHeaderColor:    color.New(color.FgGreen, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgGreen, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgGreen),
		paramColor2:         color.New(color.FgGreen),
		paramColor3:         color.New(color.FgYellow),
	})
}

func (f *loggerFactory) createLogger(l *logger) Logger {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		DisableCaller:    false,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	baseLogger, _ := cfg.Build()
	sugaredLogger := baseLogger.Sugar()
	defer sugaredLogger.Sync()
	l.SugaredLogger = sugaredLogger
	return l
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	With(args ...interface{}) Logger
	Infow(string, ...interface{})
	Named(string) Logger
	PrintGroup(msg string, args ...interface{})
	PrintSubGroup(msg string, args ...interface{})
	PrintSubGroupName(msg string, args ...interface{})
	PrintParameter(key string, value interface{})
	PrintParameters(header string, parameters map[string]interface{})
}

type logger struct {
	*zap.SugaredLogger
	groupHeaderColor    *color.Color
	subGroupHeaderColor *color.Color
	paramColor1         *color.Color
	paramColor2         *color.Color
	paramColor3         *color.Color
}

func (l *logger) With(args ...interface{}) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l *logger) Named(name string) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.Named(name),
	}
}

func (l *logger) PrintGroup(msg string, args ...interface{}) {
	l.Info(l.groupHeaderColor.Sprint(fmt.Sprintf(msg, args...)))
}

func (l *logger) PrintSubGroup(msg string, args ...interface{}) {
	l.Info("  ", fmt.Sprintf(msg, args...))
}

func (l *logger) PrintSubGroupName(msg string, args ...interface{}) {
	l.PrintSubGroup(l.subGroupHeaderColor.Sprintf(msg, args...))
}

func (l *logger) PrintParameter(key string, value interface{}) {
	l.PrintSubGroup(fmt.Sprintf("%s: %s", l.subGroupHeaderColor.Sprint(key), l.paramColor1.Sprint(value)))
}

func (l *logger) PrintParameters(header string, parameters map[string]interface{}) {
	if len(parameters) == 0 {
		return
	}

	l.PrintSubGroupName("%s:", header)
	for key, value := range parameters {
		var paramColor *color.Color
		switch key {
		case LoggerReservedKeyExpected:
			paramColor = l.paramColor2
		case LoggerReservedKeyActual:
			paramColor = l.paramColor3
		default:
			paramColor = l.paramColor1
		}

		l.Info("    ", fmt.Sprintf("%s: %s", l.subGroupHeaderColor.Sprint(key), paramColor.Sprint(value)))
	}
}
