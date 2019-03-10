package events

import (
	"errors"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/db/postgres/pg_actions"
	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/hand-writing-authentication-team/credential-store/pkg/constants"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(authReq models.AuthenticationRequest, pga *pg_actions.PgActions) error {
	username := authReq.Username
	handwriting := authReq.Handwring
	password := []byte(authReq.Password)
	encode_byte_pw, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("password encryption failed")
		return err
	}
	encoded_password := string(encode_byte_pw)

	ucm := models.UserCredentials{
		Username:        username,
		PasswordContent: encoded_password,
		Handwriting:     handwriting,
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	_, err = pga.Insert(ucm)
	if err != nil {
		log.WithError(err).Error("Internal error happen when inserting usercred")
		return err
	}
	log.Infof("successfully create user %s record!", username)
	return nil
}

func AuthUser(authReq models.AuthenticationRequest, pga *pg_actions.PgActions) error {
	username := authReq.Username
	password := []byte(authReq.Password)
	ucm, err := pga.RetrieveByUsername(username)
	if err != nil {
		return err
	}
	result := bcrypt.CompareHashAndPassword([]byte(ucm.PasswordContent), password)
	if result != nil {
		log.Debug("user password not match db's")
		return errors.New(constants.NOT_MATCH)
	}
	return nil
}

func DeleteUser(authReq models.AuthenticationRequest, pga *pg_actions.PgActions) error {
	username := authReq.Username
	err := pga.SoftDeleteByUsername(username)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(authReq models.AuthenticationRequest, pga *pg_actions.PgActions) error {
	username := authReq.Username
	handwriting := authReq.Handwring
	password := []byte(authReq.Password)
	encode_byte_pw, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("password encryption failed")
		return err
	}
	encoded_password := string(encode_byte_pw)

	ucm := models.UserCredentials{
		Username:        username,
		PasswordContent: encoded_password,
		Handwriting:     handwriting,
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	_, err = pga.Update(ucm)
	if err != nil {
		log.WithError(err).Errorf("update on user %s failed", username)
		return err
	}
	return nil
}
