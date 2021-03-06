package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const debug = true
const info = true
const warn = true
const err = true

func Debug(a ...interface{}) {
	if debug {
		log.Println(getCallerInfo("DEBUG", 2))
		log.Println(a...)
	}
}

func Info(a ...interface{}) {
	if info {
		log.Println(getCallerInfo("INFO", 2))
		log.Println(a...)
	}
}

func Warn(a ...interface{}) {
	if warn {
		log.Println(getCallerInfo("WARN", 2))
		log.Println(a...)
	}
}
func Error(a ...interface{}) {
	if err {
		log.Println(getCallerInfo("ERROR", 2))
		log.Println(a...)
	}
}

func getCallerInfo(logType string, depth int) string {
	_, file, no, ok := runtime.Caller(depth)
	if ok {
		path := file
		if strings.Contains(file, "bariot") {
			splits := strings.Split(file, "bariot")
			path = "bariot" + splits[1]
		}
		return fmt.Sprintf("[%s FROM]  %s  [LINE]  %d", logType, path, no)
	}
	return ""
}

func WithDepth(depth int, logType string, a ...interface{}) {
	if debug {
		log.Println(getCallerInfo(logType, depth))
		log.Println(a...)
	}
}
