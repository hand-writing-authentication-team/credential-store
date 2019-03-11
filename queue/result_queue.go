package queue

import (
	"encoding/json"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/models"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type ResultQueue struct {
	redisDB *redis.Client
}

func NewRedisClient(addr string) (*ResultQueue, error) {
	rq := &ResultQueue{}
	rq.redisDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	_, err := rq.redisDB.Ping().Result()
	if err != nil {
		log.WithError(err).Error("error when pinging redis, will retry in 5")
		var counter int
		for err != nil {
			counter++
			time.Sleep(5 * time.Second)
			log.Infof("retrying for the %s th time", counter)
			_, err = rq.redisDB.Ping().Result()
		}
	}
	log.Info("successfully connected to redis!")
	return rq, nil
}

func (rq *ResultQueue) SuccessInfo(authReq models.AuthenticationRequest) error {
	resultResp := models.ResultResp{
		JobID:     authReq.JobID,
		Status:    "success",
		TimeStamp: time.Now().Unix(),
	}
	var resultStr string
	resultBytes, err := json.Marshal(resultResp)
	if err != nil {
		log.WithError(err).Errorf("failed marshal the response")
		return err
	}
	resultStr = (string)(resultBytes)
	// Publish a message.
	err = rq.redisDB.Set(authReq.JobID, resultStr, 0).Err()
	if err != nil {
		log.WithError(err).Errorf("failed to give SUCCESS info for job %s", authReq.JobID)
		return err
	}
	return nil
}

func (rq *ResultQueue) FailureInfo(authReq models.AuthenticationRequest, msg string) error {
	resultResp := models.ResultResp{
		JobID:     authReq.JobID,
		Status:    "failure",
		ErrorMsg:  msg,
		TimeStamp: time.Now().Unix(),
	}
	var resultStr string
	resultBytes, err := json.Marshal(resultResp)
	if err != nil {
		log.WithError(err).Errorf("failed marshal the response")
		return err
	}
	resultStr = (string)(resultBytes)
	// Publish a message.
	err = rq.redisDB.Set(authReq.JobID, resultStr, 0).Err()
	if err != nil {
		log.WithError(err).Errorf("failed to give FAILURE info for job %s", authReq.JobID)
		return err
	}
	return nil
}

func (rq *ResultQueue) ErrorInfo(authReq models.AuthenticationRequest, msg string) error {
	resultResp := models.ResultResp{
		JobID:     authReq.JobID,
		Status:    "error",
		ErrorMsg:  msg,
		TimeStamp: time.Now().Unix(),
	}
	var resultStr string
	resultBytes, err := json.Marshal(resultResp)
	if err != nil {
		log.WithError(err).Errorf("failed marshal the response")
		return err
	}
	resultStr = (string)(resultBytes)
	// Publish a message.
	err = rq.redisDB.Set(authReq.JobID, resultStr, 0).Err()
	if err != nil {
		log.WithError(err).Errorf("failed to give ERROR info for job %s", authReq.JobID)
		return err
	}
	return nil
}
