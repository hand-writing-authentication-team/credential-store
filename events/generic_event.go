package events

import (
	"encoding/json"

	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/hand-writing-authentication-team/credential-store/pkg/constants"
	"github.com/hand-writing-authentication-team/credential-store/queue"
	log "github.com/sirupsen/logrus"
)

func GenericEventHandler(msg []byte, queueClient *queue.Queue) (err error) {
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
		break
	case constants.AuthAction:
		log.Info("start to authenticate user")
		break
	default:
		break
	}
	return nil
}
