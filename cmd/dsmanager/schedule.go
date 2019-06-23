package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	A scheduled task can be setup and monitored
 */

var scheduled = make(map[string]*scheduledTask)
var scheduledMutex = sync.Mutex{}

type scheduledTask struct {
	interval int //time in minutes between runs
	timeout time.Duration

	Running bool

	LastRun   time.Time
	LastError error

	task func() error
}

func init(){
	go startScheduler()
}

func startScheduler() {
	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		scheduledMutex.Lock()

		for k, v := range scheduled {
			if time.Now().After(v.LastRun.Add(time.Minute * time.Duration(v.interval))) {
				go runTask(k, v)
			}
		}

		scheduledMutex.Unlock()
	}
}

func runTask(label string, sc *scheduledTask) {
	if sc.Running {
		webLogf("Scheduler: Cannot start %v: already Running", label)
		return
	}
	webLogf("Scheduler: Starting %v", label)
	sc.LastRun = time.Now()
	timer := time.NewTimer(sc.timeout)
	context := make(chan error)

	sc.Running = true
	defer func(){
		sc.Running = false
	}()
	go func() {
		select {
		case context<- sc.task():
			return
		case  <-context: //Task has been timed out :(
			return
		}
	}()

	select {
	case <-timer.C: //timedout
		context <- fmt.Errorf("timedout")
		webLogf("Scheduler: %v reached timeout", label)
		sc.LastError = fmt.Errorf("task reached timeout")
		return
	case err := <-context: //returned
		if err != nil {
			webLogf("Scheduler: %v returned an error %v", label, err)
			sc.LastError = err
		} else {
			webLogf("Scheduler: Completed %v", label)
			sc.LastError = nil
		}
		return
	}
}

func start(sc scheduledTask, label string){
	scheduledMutex.Lock()
	sc.LastRun = time.Now()
	scheduled[label] = &sc
	scheduledMutex.Unlock()
}

func stop(label string){
	scheduledMutex.Lock()
	delete(scheduled, label)
	scheduledMutex.Unlock()
}