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

func GenericEventHandler(msg []byte, queueClient *queue.Queue, pga *pg_actions.PgActions) (err error) {
	var authReq models.AuthenticationRequest
	err = json.Unmarshal(msg, &authReq)
	if err != nil {
		// TODO: throw into a error exchange
		log.WithField("error", err).Error("cannot unmarshal the message")
		return err
	}

	switch authReq.Action {
	case constants.CreateAction:
		// TODO: when postgres dao is done, add this logic in
		log.Info("start to create account")
		err = CreateUser(authReq, pga)
		if err != nil {
			// TODO: should go into the result queue
			return err
		}
		break
	case constants.AuthAction:
		log.Info("start to authenticate user")
		err = AuthUser(authReq, pga)
		if err != nil {
			switch err.Error() {
			case sql.ErrNoRows.Error():
				// TODO: give back and tell no record
				log.Error("There is no record for this auth request")
				return err
			case constants.NOT_MATCH:
				// TODO: auth failed because not match
				log.Error("The password does not match what is in DB")
				return err
			default:
				// TODO: internal server error
				log.WithError(err).Error("Internal server error")
				return err
			}
		}
		log.Infof("user %s has authenticate password successfully", authReq.Username)
		break
	case constants.UpdateAction:
		break
	case constants.DeleteAction:
		break
	default:
		break
	}
	return nil
}
