package logger

import (
	"fmt"
	tm "github.com/buger/goterm"
)

type Logger struct {
}

func (l *Logger) PrintH1(msg string, args ...interface{}) {
	tm.Println(tm.Color(tm.Bold(fmt.Sprintf(msg, args...)), tm.GREEN))
	tm.Flush()
}

func (l *Logger) PrintH2(msg string, args ...interface{}) {
	//tm.Println(tm.MoveTo(tm.Color(tm.Bold(fmt.Sprintf(msg, args...)), tm.YELLOW), 10|tm.PCT, tm.CurrentHeight()))
	tm.Println(tm.Color(tm.Bold(fmt.Sprintf(msg, args...)), tm.YELLOW))
	tm.Flush()
}

func (l *Logger) PrintH2Parameters(header string, parameters map[string]interface{}) {
	if len(parameters) == 0 {
		return
	}

	tm.Println(tm.Color(tm.Bold(fmt.Sprintf("%s:", header)), tm.YELLOW))
	for key, value := range parameters {
		tm.Println(tm.Color(fmt.Sprintf("%s:%v", key, value), tm.CYAN))
	}
	tm.Flush()
}
