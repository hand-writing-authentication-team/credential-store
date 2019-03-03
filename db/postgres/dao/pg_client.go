package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/hand-writing-authentication-team/credential-store/models"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PGDBInstance struct {
	conn *sql.DB
}

func scanUserCredRow(row *sql.Row, ucm *models.UserCredentials) error {
	return row.Scan(&ucm.ID,
		&ucm.Username,
		&ucm.Handwriting,
		&ucm.PasswordContent,
		&ucm.Created,
		&ucm.Modified,
		&ucm.Deleted)
}

func NewDBInstance(dbhost, dbport, dbuser, dbpass, dbname string) (*PGDBInstance, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.WithField("error", err).Error("cannot open postgres server")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.WithField("error", err).Error("Ping to postgres failed")
		return nil, err
	}
	pgInstance := PGDBInstance{
		conn: db,
	}
	log.Infof("Successfully connected!")
	return &pgInstance, nil
}

func (p *PGDBInstance) Begin() (txm *TXManager, err error) {
	if p == nil {
		log.Error("connection is null")
		return nil, errors.New("connection is null")
	}
	tx, err := p.conn.Begin()
	if err != nil {
		log.WithError(err).Error("Error happen when start new transaction")
		return nil, err
	}
	txm = &TXManager{}
	txm.Tx = tx
	return txm, nil
}

func (p *PGDBInstance) Insert(txm *TXManager, ucm models.UserCredentials) (*models.UserCredentials, error) {
	stmt := `INSERT INTO user_cred (username, hand_writing, pw_encoded, created, modified, deleted)
                     VALUES($1,$2,$3,$4,$5,FALSE) RETURNING *;`

	row := txm.Tx.QueryRow(stmt, ucm.Username, ucm.Handwriting,
		ucm.PasswordContent, ucm.Created, ucm.Modified)
	retUcm := &models.UserCredentials{}
	err := scanUserCredRow(row, retUcm)
	defer txm.End(err)
	if err != nil {
		log.WithError(err).Errorf("insertion for user %s failed", ucm.Username)
		return nil, err
	}
	log.Infof("Insertion for user %s succeeded!", ucm.Username)
	return retUcm, nil
}
