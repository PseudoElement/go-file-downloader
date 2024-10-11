package logger

import (
	"log"
	"time"
)

type LogMsg struct {
	msg       string
	timestamp int64
}

type Logger struct {
	logsByFuncName map[string][]LogMsg
}

func NewLogger() Logger {
	return Logger{
		logsByFuncName: make(map[string][]LogMsg),
	}
}

func (l *Logger) AddLog(funcName string, msg string) {
	_, ok := l.logsByFuncName[funcName]
	if !ok {
		l.logsByFuncName[funcName] = make([]LogMsg, 0, 20)
	}

	msgWithTime := LogMsg{
		msg:       msg,
		timestamp: time.Now().UnixMilli(),
	}

	l.logsByFuncName[funcName] = append(l.logsByFuncName[funcName], msgWithTime)
}

func (l *Logger) ShowLogs(funcName string) {
	if logs, ok := l.logsByFuncName[funcName]; !ok {
		log.Printf("[ShowLogs] Invalid funcName %v\n", funcName)
		return
	} else {
		log.Printf("===================START_%v===================\n", funcName)
		for ind, el := range logs {
			if ind == 0 {
				log.Printf("First log. %v\n", el.msg)
			} else {
				prevEl := logs[ind-1]
				msFromPrevLog := el.timestamp - prevEl.timestamp
				log.Printf("After %v ms. %v\n", msFromPrevLog, el.msg)
			}
		}

		totalTime := time.Now().UnixMilli() - logs[0].timestamp

		log.Printf("Total time - %v ms.\n", totalTime)
		log.Printf("===================END_%v===================\n\n\n", funcName)

		l.clearLogsByFuncName(funcName)
	}
}

func (l *Logger) clearLogsByFuncName(funcName string) {
	l.logsByFuncName[funcName] = nil
}
