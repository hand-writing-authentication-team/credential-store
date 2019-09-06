package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/hand-writing-authentication-team/credential-store/pkg/constants"
	"github.com/lib/pq"
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
		&ucm.UserModel,
		&ucm.PasswordContent,
		&ucm.Race,
		&ucm.Created,
		&ucm.Modified,
		&ucm.Deleted)
}

func scanUserHWRow(row *sql.Row, uvHW *models.UserValidateHW) error {
	return row.Scan(&uvHW.ID,
		&uvHW.UserID,
		&uvHW.Username,
		&uvHW.Handwriting,
		&uvHW.Created,
		&uvHW.Modified,
		&uvHW.Deleted)
}

func NewDBInstance(dbhost, dbport, dbuser, dbpass, dbname string) (*PGDBInstance, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.WithField("error", err).Error("cannot open postgres server")
		log.Info("Will start retrying")
		counter := 0
		for err != nil {
			counter++
			time.Sleep(5 * time.Second)
			log.Infof("restart %s th times", counter)
			db, err = sql.Open("postgres", psqlInfo)
		}
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
	stmt := `INSERT INTO user_cred (username, hand_writing, user_model, pw_encoded, race, created, modified, deleted)
                     VALUES($1,$2,$3,$4,$5,$6,$7,FALSE) RETURNING *;`

	row := txm.Tx.QueryRow(stmt, ucm.Username, ucm.Handwriting, ucm.UserModel,
		ucm.PasswordContent, ucm.Race, ucm.Created, ucm.Modified)
	retUcm := &models.UserCredentials{}
	err := scanUserCredRow(row, retUcm)
	if err != nil {
		log.WithError(err).Errorf("insertion for user %s failed", ucm.Username)
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			return nil, errors.New(constants.ACCOUNT_ALREADY_EXIST)
		}
		return nil, err
	}
	log.Infof("Insertion for user %s succeeded!", ucm.Username)
	return retUcm, nil
}

func (p *PGDBInstance) RetrieveByUsername(txm *TXManager, username string) (*models.UserCredentials, error) {
	stmt := `SELECT * FROM user_cred WHERE username = $1 AND deleted = FALSE;`

	row := txm.Tx.QueryRow(stmt, username)
	retUcm := &models.UserCredentials{}
	err := scanUserCredRow(row, retUcm)
	if err != nil {
		log.WithError(err).Errorf("retrieval for user %s failed", username)
		return nil, err
	}
	log.Infof("retrieval for user %s succeeded!", username)
	return retUcm, nil
}

func (p *PGDBInstance) SoftDeleteByUsername(txm *TXManager, username string) error {
	stmt := `UPDATE user_cred SET deleted = TRUE WHERE username = $1 AND deleted = FALSE;`

	_, err := txm.Tx.Exec(stmt, username)
	if err != nil {
		log.WithError(err).Errorf("soft delete %s failed", username)
		return err
	}
	log.Infof("soft delete for user %s succeeded!", username)
	return nil
}

func (p *PGDBInstance) Update(txm *TXManager, ucm models.UserCredentials) (*models.UserCredentials, error) {
	stmt := `UPDATE user_cred 
			SET hand_writing = $1,
			pw_encoded = $2,
			created = $3, 
			modified = $4 WHERE username = $5 AND deleted = FALSE
			RETURNING *;`

	row := txm.Tx.QueryRow(stmt, ucm.Handwriting,
		ucm.PasswordContent, ucm.Created, ucm.Modified, ucm.Username)
	retUcm := &models.UserCredentials{}
	err := scanUserCredRow(row, retUcm)
	if err != nil {
		log.WithError(err).Errorf("updation for user %s failed", ucm.Username)
		return nil, err
	}
	log.Infof("updation for user %s succeeded!", ucm.Username)
	return retUcm, nil
}

func (p *PGDBInstance) InsertSecondHandwriting(txm *TXManager, uvHW models.UserValidateHW) (*models.UserValidateHW, error) {
	ucm, err := p.RetrieveByUsername(txm, uvHW.Username)
	if err != nil {
		log.WithError(err).Errorf("[SEC_HW]retrieve user cred model failed for user %s", uvHW.Username)
		return nil, errors.New(constants.ACCOUNT_NOT_FOUND)
	}
	userID := ucm.ID

	stmt := `INSERT INTO validate_handwriting (userId, username, hand_writing, created, modified, deleted)
                     VALUES($1,$2,$3,$4,$5,FALSE) RETURNING *;`

	row := txm.Tx.QueryRow(stmt, userID, uvHW.Username, uvHW.Handwriting,
		uvHW.Created, uvHW.Modified)
	retUvHW := &models.UserValidateHW{}
	err = scanUserHWRow(row, retUvHW)
	if err != nil {
		log.WithError(err).Errorf("[SEC_HW]insertion for user's second handwriting %s failed", uvHW.Username)
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			return nil, errors.New(constants.ACCOUNT_ALREADY_EXIST)
		}
		return nil, err
	}
	log.Infof("[SEC_HW]Insertion for user's second handwriting %s succeeded!", uvHW.Username)
	return retUvHW, nil
}
