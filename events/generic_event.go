package events

import (
	"database/sql"
	"encoding/json"

	"github.com/hand-writing-authentication-team/credential-store/db/postgres/pg_actions"
	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/hand-writing-authentication-team/credential-store/pkg/constants"
	"github.com/hand-writing-authentication-team/credential-store/queue"
	log "github.com/sirupsen/logrus"
)

func GenericEventHandler(msg []byte, queueClient *queue.Queue, pga *pg_actions.PgActions, rq *queue.ResultQueue) {
	go func() {
		log.Info("event thread start")
		var authReq models.AuthenticationRequest
		err := json.Unmarshal(msg, &authReq)
		if err != nil {
			// TODO: throw into a error exchange
			log.WithField("error", err).Error("cannot unmarshal the message")
			return
		}

		switch authReq.Action {
		case constants.CreateAction:
			// TODO: when postgres dao is done, add this logic in
			log.Info("start to create account")
			err = CreateUser(authReq, pga)
			if err != nil {
				// TODO: should go into the result queue
				log.WithError(err).Error("Error occured when creating user")
				rq.ErrorInfo(authReq, err.Error())
				return
			}
			rq.SuccessInfo(authReq)
			break
		case constants.AuthAction:
			log.Info("start to authenticate user")
			err = AuthUser(authReq, pga)
			if err != nil {
				switch err.Error() {
				case sql.ErrNoRows.Error():
					rq.FailureInfo(authReq, err.Error())
					log.Error("There is no record for this auth request")
					return
				case constants.NOT_MATCH:
					rq.FailureInfo(authReq, err.Error())
					log.Error("The password does not match what is in DB")
					return
				default:
					rq.ErrorInfo(authReq, err.Error())
					log.WithError(err).Error("Internal server error")
					return
				}
			}
			log.Infof("user %s has authenticate password successfully", authReq.Username)
			rq.SuccessInfo(authReq)
			break
		case constants.UpdateAction:
			log.Infof("start to update user")
			err = UpdateUser(authReq, pga)
			if err != nil {
				switch err.Error() {
				case sql.ErrNoRows.Error():
					rq.FailureInfo(authReq, err.Error())
					log.Error("There is no record for this update request")
					return
				default:
					rq.ErrorInfo(authReq, err.Error())
					log.WithError(err).Error("Internal server error")
					return
				}
			}
			rq.SuccessInfo(authReq)
			break
		case constants.DeleteAction:
			log.Infof("start to delete user")
			err = DeleteUser(authReq, pga)
			if err != nil {
				switch err.Error() {
				case sql.ErrNoRows.Error():
					rq.FailureInfo(authReq, err.Error())
					log.Error("There is no record for this delete request")
					return
				default:
					rq.ErrorInfo(authReq, err.Error())
					log.WithError(err).Error("Internal server error")
					return
				}
			}
			rq.SuccessInfo(authReq)
			break
		default:
			break
		}
		log.Info("Successfully complete the event")
		return
	}()
}
