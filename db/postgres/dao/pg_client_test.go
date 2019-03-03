package dao

import (
	"os"
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
	}
	defer txm.Rollback()
	retUcm, err := PGInstance.Insert(txm, ucm)
	if err != nil {
		t.Errorf("met an error %s when inserting", err)
	}
	assert.Equal(t, ucm.Username, retUcm.Username)
}
