package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var webLogLines = make([]string, 0)
var webLogMutex = sync.Mutex{}
var webLogTemplate = template.Must(template.New("base").ParseFiles("templates/weblogs.html"))

type webLogsPage struct {
	Lines []string
}

func handleWebLogs(res http.ResponseWriter, req *http.Request) {
	webLogMutex.Lock()
	err := webLogTemplate.Execute(res, webLogsPage{Lines: webLogLines})
	webLogMutex.Unlock()
	if err != nil {
		http.Error(res, fmt.Sprintf("template execution failed: %v", err), http.StatusInternalServerError)
		return
	}
}

func webLogf(format string, v ...interface{}) {
	line := fmt.Sprintf(format, v...)
	log.Println(line)
	webLogMutex.Lock()
	webLogLines = append(webLogLines, fmt.Sprintf("%v: %v", time.Now().Format(time.Stamp), line))
	webLogMutex.Unlock()
}

func webLogCleanup() scheduledTask {
	sc := scheduledTask{
		interval: 60,
		timeout:  time.Minute,

		task: func() error {
			if len(webLogLines) > 1000 {
				webLogMutex.Lock()
				webLogLines = webLogLines[len(webLogLines)-1000:]
				webLogMutex.Unlock()
			}
			return nil
		},
	}
	return sc
}