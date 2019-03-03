package pg_actions

import (
	"github.com/hand-writing-authentication-team/credential-store/db/postgres/dao"
	"github.com/hand-writing-authentication-team/credential-store/models"
	log "github.com/sirupsen/logrus"
)

type PgActions struct {
	db *dao.PGDBInstance
}

func NewPgActions(dbconn *dao.PGDBInstance) *PgActions {
	pga := &PgActions{
		db: dbconn,
	}
	return pga
}

func (p *PgActions) Insert(ucm models.UserCredentials) (*models.UserCredentials, error) {
	txm, err := p.db.Begin()
	if err != nil {
		log.Error("transaction creation failed")
		return nil, err
	}
	retUcm, err := p.db.Insert(txm, ucm)
	defer txm.End(err)
	if err != nil {
		log.WithError(err).Error("error occured when insert a user cred model")
		return nil, err
	}
	log.Debug("successfully inserted usercred model")
	return retUcm, nil
}

func (p *PgActions) RetrieveByUsername(username string) (*models.UserCredentials, error) {
	txm, err := p.db.Begin()
	if err != nil {
		log.Error("transaction creation failed")
		return nil, err
	}
	retUcm, err := p.db.RetrieveByUsername(txm, username)
	defer txm.End(err)
	if err != nil {
		log.WithError(err).Error("error occured when retrieve a user cred model")
		return nil, err
	}
	log.Debug("successfully retrieval usercred model")
	return retUcm, nil
}

func (p *PgActions) SoftDeleteByUsername(username string) error {
	txm, err := p.db.Begin()
	if err != nil {
		log.Error("transaction creation failed")
		return err
	}
	err = p.db.SoftDeleteByUsername(txm, username)
	defer txm.End(err)
	if err != nil {
		log.WithError(err).Error("error occured when delete a user cred model")
		return err
	}
	log.Debug("successfully deleted usercred model")
	return nil
}

func (p *PgActions) Update(ucm models.UserCredentials) (*models.UserCredentials, error) {
	txm, err := p.db.Begin()
	if err != nil {
		log.Error("transaction creation failed")
		return nil, err
	}
	retUcm, err := p.db.Update(txm, ucm)
	defer txm.End(err)
	if err != nil {
		log.WithError(err).Error("error occured when updating a user cred model")
		return nil, err
	}
	log.Debug("successfully updated usercred model")
	return retUcm, nil
}
