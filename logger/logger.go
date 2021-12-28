package logger

import (
	"fmt"
	tm "github.com/buger/goterm"
)

type colors struct {
	h1    int
	h2    int
	param int
}

var (
	defaultColors = colors{
		h1:    tm.YELLOW,
		h2:    tm.WHITE,
		param: tm.WHITE,
	}
	errorColors = colors{
		h1:    tm.YELLOW,
		h2:    tm.RED,
		param: tm.RED,
	}
)

type Logger interface {
	PrintGroup(msg string, args ...interface{})
	PrintSubGroup(msg string, args ...interface{})
	PrintParameters(header string, parameters map[string]interface{})
	PrintParameter(key string, value interface{})
	PrintFailure()
	PrintSuccess()
	ToErrorColors() Logger
}

func NewFxLogger() Logger {
	return &stdoutLogger{
		colors: defaultColors,
	}
}

type stdoutLogger struct {
	colors colors
}

func (l *stdoutLogger) PrintGroup(msg string, args ...interface{}) {
	tm.Println(tm.Color(tm.Bold(fmt.Sprintf(msg, args...)), l.colors.h1))
	tm.Flush()
}

func (l *stdoutLogger) PrintSubGroup(msg string, args ...interface{}) {
	tm.MoveCursorForward(2)
	tm.Println(tm.Color(tm.Bold(fmt.Sprintf(msg, args...)), l.colors.h2))
	tm.Flush()
}

func (l *stdoutLogger) PrintParameters(header string, parameters map[string]interface{}) {
	if len(parameters) == 0 {
		return
	}

	l.PrintSubGroup(fmt.Sprintf("%s:", header))
	for key, value := range parameters {
		l.PrintParameter(key, value)
	}
}

func (l *stdoutLogger) PrintParameter(key string, value interface{}) {
	tm.MoveCursorForward(4)
	tm.Println(tm.Color(fmt.Sprintf("%s: %v", key, value), l.colors.param))
	tm.Flush()
}

func (l *stdoutLogger) PrintFailure() {
	tm.Println(tm.Color(tm.Bold("Status: Failure"), tm.RED))
	tm.Flush()
}

func (l *stdoutLogger) PrintSuccess() {
	tm.Println(tm.Color(tm.Bold("Status: Success"), tm.GREEN))
	tm.Flush()
}

func (l *stdoutLogger) ToErrorColors() Logger {
	return &stdoutLogger{
		colors: errorColors,
	}
}
