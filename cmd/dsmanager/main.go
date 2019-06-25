package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var startTime = time.Now()
var indexTemplate = template.Must(template.New("base").ParseFiles("templates/index.html"))

func main(){
	//Schedule some shit
	start(webLogCleanup(), "WebLogCleanup")
	start(plexHibernate(), "PlexHibernate")

	//Setup webserver
	router := http.NewServeMux()

	router.HandleFunc("/weblogs", handleWebLogs)
	router.HandleFunc("/hibernate", handleHibernate)
	router.HandleFunc("/", handleIndex)

	webLogf("Starting HTTP server on port :8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}

type indexPage struct {
	ScheduledTasks map[string]*scheduledTask
	StartTime time.Time
	NextOnTime time.Time
	NextOffTime time.Time
}

func handleIndex(res http.ResponseWriter, req *http.Request) {

	ip := indexPage{ScheduledTasks:scheduled, StartTime:startTime, NextOnTime: getNextOn(), NextOffTime: getNextOff()}

	err := indexTemplate.Execute(res, ip)
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleHibernate(res http.ResponseWriter, req *http.Request) {
	err := hibernate()
	if err != nil {
		webLogf("HandleHibernate: hibernate() failed: %v", err)
		http.Error(res, fmt.Sprintf("hibernate() failed: %v", err), http.StatusInternalServerError)
		return
	}

	res.Write([]byte("Hibernating Server - dsmanager will become unavailable shortly"))
}

func shutDownTill(till time.Time) error {
	d := time.Until(till)

	cmd := exec.Command("rtcwake", "-s", fmt.Sprintf("%d", int(d.Seconds())), "-m", "off")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rtcwake cmd failed: \nRTCWAKE OUTPUT:\n%s", string(b))
	}

	return nil
}