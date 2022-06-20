package debug

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const debug = true

func LogDebug(a ...interface{}) {
	if debug {
		log.Println(getCallerInfo("DEBUG", 2))
		log.Println(a...)
	}
}

func LogInfo(a ...interface{}) {
	if debug {
		log.Println(getCallerInfo("INFO", 2))
		log.Println(a...)
	}
}

func LogWarn(a ...interface{}) {
	if debug {
		log.Println(getCallerInfo("WARN", 2))
		log.Println(a...)
	}
}
func LogError(a ...interface{}) {
	if debug {
		log.Println(getCallerInfo("ERROR", 2))
		log.Println(a...)
	}
}

func getCallerInfo(logType string, depth int) string {
	_, file, no, ok := runtime.Caller(depth)
	if ok {
		splits := strings.Split(file, "/")
		fileName := splits[len(splits)-1]
		return fmt.Sprintf("[%s FROM]  %s  [LINE]  %d", logType, fileName, no)
	}
	return ""
}

func LogWithDepth(depth int, logType string, a ...interface{}) {
	if debug {
		log.Println(getCallerInfo(logType, depth))
		log.Println(a...)
	}
}
