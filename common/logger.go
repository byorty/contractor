package common

import (
	"fmt"
	"github.com/fatih/color"
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
	return &logger{
		groupHeaderColor:    color.New(color.FgBlue, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgMagenta, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgWhite),
	}
}

func (f *loggerFactory) CreateErrorLogger() Logger {
	return &logger{
		groupHeaderColor:    color.New(color.FgRed, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgRed, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgRed),
		paramColor2:         color.New(color.FgGreen),
		paramColor3:         color.New(color.FgYellow),
	}
}

func (f *loggerFactory) CreateSuccessLogger() Logger {
	return &logger{
		groupHeaderColor:    color.New(color.FgGreen, color.Bold, color.Underline),
		subGroupHeaderColor: color.New(color.FgGreen, color.Bold, color.Underline),
		paramColor1:         color.New(color.FgGreen),
	}
}

type Logger interface {
	PrintGroup(msg string, args ...interface{})
	PrintSubGroup(msg string, args ...interface{})
	PrintSubGroupName(msg string, args ...interface{})
	PrintParameter(key string, value interface{})
	PrintParameters(header string, parameters map[string]interface{})
}

type logger struct {
	groupHeaderColor    *color.Color
	subGroupHeaderColor *color.Color
	paramColor1         *color.Color
	paramColor2         *color.Color
	paramColor3         *color.Color
}

func (l *logger) PrintGroup(msg string, args ...interface{}) {
	l.groupHeaderColor.Println(fmt.Sprintf(msg, args...))
}

func (l *logger) PrintSubGroup(msg string, args ...interface{}) {
	fmt.Println("  ", fmt.Sprintf(msg, args...))
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

		fmt.Println("    ", fmt.Sprintf("%s: %s", l.subGroupHeaderColor.Sprint(key), paramColor.Sprint(value)))
	}
}
