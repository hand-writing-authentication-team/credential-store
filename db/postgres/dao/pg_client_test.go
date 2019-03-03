package dao

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hand-writing-authentication-team/credential-store/models"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/gommon/log"
)

func setup() *PGDBInstance {
	pgHost := strings.TrimSpace(os.Getenv("PG_HOST"))
	pgUser := strings.TrimSpace(os.Getenv("PG_USER"))
	pgPassword := strings.TrimSpace(os.Getenv("PG_PASSWORD"))
	pgPort := strings.TrimSpace(os.Getenv("PG_PORT"))
	pgDB := strings.TrimSpace(os.Getenv("PG_DBNAME"))

	PGInstance, err := NewDBInstance(pgHost, pgPort, pgUser, pgPassword, pgDB)
	if err != nil {
		log.Error("pg connection failed, test cannot proceed")
		os.Exit(1)
	}
	log.Info("start to connect to postgres")
	PGInstance.conn.Ping()
	return PGInstance
}

func TestInsertion(t *testing.T) {
	PGInstance := setup()
	ucm := models.UserCredentials{
		Username:        "todd-test",
		PasswordContent: "abcd1234",
		Handwriting:     "A1B2C3D4",
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	txm, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
		return
	}
	defer txm.Rollback()
	retUcm, err := PGInstance.Insert(txm, ucm)
	if err != nil {
		t.Errorf("met an error %s when inserting", err)
		return
	}
	assert.Equal(t, ucm.Username, retUcm.Username)
}

func TestSoftDeleteByUsername(t *testing.T) {
	PGInstance := setup()
	ucm := models.UserCredentials{
		Username:        "todd-test",
		PasswordContent: "abcd1234",
		Handwriting:     "A1B2C3D4",
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	txm, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	defer txm.Rollback()
	_ = insertUser(PGInstance, txm, ucm)
	txm1, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	err = PGInstance.SoftDeleteByUsername(txm1, "todd-test")
	if err != nil {
		t.Errorf("met error %s when deletion", err)
		return
	}
	txm1.End(nil)
}

func TestRetrivalByUsername(t *testing.T) {
	PGInstance := setup()
	ucm := models.UserCredentials{
		Username:        "todd-test",
		PasswordContent: "abcd1234",
		Handwriting:     "A1B2C3D4",
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	txm, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	expUcm := insertUser(PGInstance, txm, ucm)
	txm.End(nil)
	txm1, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	retUcm, err := PGInstance.RetrieveByUsername(txm1, "todd-test")
	if err != nil {
		t.Errorf("met error %s when retrieval", err)
		return
	}
	txm1.End(nil)
	assert.True(t, reflect.DeepEqual(expUcm, retUcm))
	softdeleteUser(PGInstance, "todd-test")
}

func TestUpdateByUsername(t *testing.T) {
	PGInstance := setup()
	ucm := models.UserCredentials{
		Username:        "todd-test",
		PasswordContent: "abcd1234",
		Handwriting:     "A1B2C3D4",
		Created:         time.Now().Unix(),
		Modified:        time.Now().Unix(),
	}
	txm, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	expUcm := insertUser(PGInstance, txm, ucm)
	txm.End(nil)
	expUcm.PasswordContent = "1234abcd"
	txm1, err := PGInstance.Begin()
	if err != nil {
		t.Error("transaction creation failed")
	}
	retUcm, err := PGInstance.Update(txm1, *expUcm)
	if err != nil {
		t.Errorf("met error %s when retrieval", err)
		return
	}
	txm1.End(nil)
	assert.Equal(t, retUcm.PasswordContent, expUcm.PasswordContent)
	softdeleteUser(PGInstance, "todd-test")
}

func insertUser(PGInstance *PGDBInstance, txm *TXManager, ucm models.UserCredentials) *models.UserCredentials {
	retUcm, err := PGInstance.Insert(txm, ucm)
	if err != nil {
		panic(err)
	}
	return retUcm
}

func softdeleteUser(PGInstance *PGDBInstance, username string) {
	txm, err := PGInstance.Begin()
	if err != nil {
		panic(err)
	}
	err = PGInstance.SoftDeleteByUsername(txm, username)
	if err != nil {
		panic(err)
	}
	txm.End(err)
	return
}
