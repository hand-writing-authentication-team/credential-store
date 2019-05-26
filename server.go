package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/clients"
	"github.com/hand-writing-authentication-team/credential-store/db/postgres/pg_actions"

	"github.com/hand-writing-authentication-team/credential-store/db/postgres/dao"
	"github.com/hand-writing-authentication-team/credential-store/events"
	"github.com/hand-writing-authentication-team/credential-store/queue"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type config struct {
	xizhiEnabled string
	xizhiUrl     string

	mqHost     string
	mqPort     string
	mqUsername string
	mqPassword string

	PgAct    *pg_actions.PgActions
	xzClient *clients.XizhiClient
	QC       *queue.Queue
	RQ       *queue.ResultQueue
	ch       <-chan amqp.Delivery
}

var goroutineDelta = make(chan int)
var serverConfig config

func start() {
	serverConfig.xizhiEnabled = strings.TrimSpace(os.Getenv("XIZHI_ENABLED"))
	if serverConfig.xizhiEnabled == "true" {
		log.Info("xizhi server is enabled, time to check how you write")
		serverConfig.xizhiUrl = strings.TrimSpace("XIZHI_URL")
		xzc, err := clients.NewXizhiClient(serverConfig.xizhiUrl, time.Duration(5*time.Second))
		if err != nil {
			log.Errorf("met error when creating xz client, will bootstrap as it disabled")
			serverConfig.xizhiEnabled = "false"
		} else {
			serverConfig.xzClient = xzc
			log.Info("xizhi client configured successfully")
		}
	} else {
		log.Info("xizhi server is not enabled, only checking password")
	}

	serverConfig.mqHost = os.Getenv("MQ_HOST")
	serverConfig.mqPort = os.Getenv("MQ_PORT")
	serverConfig.mqUsername = os.Getenv("MQ_USER")
	serverConfig.mqPassword = os.Getenv("MQ_PASSWORD")

	pgHost := strings.TrimSpace(os.Getenv("PG_HOST"))
	pgUser := strings.TrimSpace(os.Getenv("PG_USER"))
	pgPassword := strings.TrimSpace(os.Getenv("PG_PASSWORD"))
	pgPort := strings.TrimSpace(os.Getenv("PG_PORT"))
	pgDB := strings.TrimSpace(os.Getenv("PG_DBNAME"))

	redisAddr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))

	if strings.TrimSpace(serverConfig.mqHost) == "" || strings.TrimSpace(serverConfig.mqPassword) == "" || strings.TrimSpace(serverConfig.mqPort) == "" || strings.TrimSpace(serverConfig.mqUsername) == "" {
		log.Fatal("one of the mq config env is not set!")
		os.Exit(1)
	}

	if pgHost == "" || pgUser == "" || pgPassword == "" || pgPort == "" || pgDB == "" {
		log.Fatal("one of the postgres configuration is not set")
		os.Exit(1)
	}

	if redisAddr == "" {
		log.Fatal("one of the redis configuration is not set")
		os.Exit(1)
	}

	queueClient, err := queue.NewQueueInstance(serverConfig.mqHost, serverConfig.mqPort, serverConfig.mqUsername, serverConfig.mqPassword)
	if err != nil {
		os.Exit(1)
	}
	serverConfig.QC = queueClient

	pgConn, err := dao.NewDBInstance(pgHost, pgPort, pgUser, pgPassword, pgDB)
	if err != nil {
		os.Exit(1)
	}
	serverConfig.PgAct = pg_actions.NewPgActions(pgConn)

	serverConfig.RQ, err = queue.NewRedisClient(redisAddr)
	if err != nil {
		os.Exit(1)
	}

	events.XZClient = serverConfig.xzClient
	log.Infof("everything is in place")
	return
}

func main() {

	log.Info("start to bootstrap credential-store server")
	start()
	gentlyExit()
	// set up listener logic
	var err error
	serverConfig.ch, err = serverConfig.QC.Consume("credstoreIn")
	if err != nil {
		log.Fatal("queue is not declared")
		os.Exit(1)
	}

	go forever()

	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		if numGoroutines == 0 {
			log.Info("no more requests, waiting")
		}
	}
}

// Conceptual code
func forever() {
	log.Info("server started.")
	foreverRunner := make(chan bool)
	go func() {
		for d := range serverConfig.ch {
			goroutineDelta <- 1
			events.GenericEventHandler(d.Body, serverConfig.QC, serverConfig.PgAct, serverConfig.RQ)
			goroutineDelta <- -1
		}
	}()
	log.Info("waiting...")
	<-foreverRunner
}

func gentlyExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := serverConfig.QC.DestroyQueueInstance()
		if err != nil {
			log.WithField("error", err).Fatal("queue connection close failed")
			os.Exit(1)
		} else {
			log.Info("queue connection closed properly, gently quitting")
			os.Exit(0)
		}
	}()
}
