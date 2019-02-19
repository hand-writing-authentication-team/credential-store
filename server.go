package main

import (
	"log"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/hand-writing-authentication-team/credential-store/queue"
)

type config struct {
	mqHost     string
	mqPort     string
	mqUsername string
	mqPassword string

	QC *queue.Queue
}

var goroutineDelta = make(chan int)
var serverConfig config

func start() {
	serverConfig.mqHost = os.Getenv("MQ_HOST")
	serverConfig.mqPort = os.Getenv("MQ_PORT")
	serverConfig.mqUsername = os.Getenv("MQ_USER")
	serverConfig.mqPassword = os.Getenv("MQ_PASSWORD")

	if strings.TrimSpace(serverConfig.mqHost) == "" || strings.TrimSpace(serverConfig.mqPassword) == "" || strings.TrimSpace(serverConfig.mqPort) == "" || strings.TrimSpace(serverConfig.mqUsername) == "" {
		log.Fatal("one of the mq config env is not set!")
		os.Exit(1)
	}

	queueClient, err := queue.NewQueueInstance(serverConfig.mqHost, serverConfig.mqPort, serverConfig.mqUsername, serverConfig.mqPassword)
	if err != nil {
		os.Exit(1)
	}
	serverConfig.QC = queueClient
	return
}

func main() {
	glog.Info("start to bootstrap credential-store server")
	start()
	go forever()

	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		if numGoroutines == 0 {
			os.Exit(0)
		}
	}
}

// Conceptual code
func forever() {
	for {
		if needToCreateANewGoroutine {
			// Make sure to do this before "go f()", not within f()
			goroutineDelta <- +1

			go f()
		}
	}
}

func f() {
	// When the termination condition for this goroutine is detected, do:
	goroutineDelta <- -1
}
