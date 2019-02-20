package dao

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func NewDBInstance(dbhost, dbport, dbuser, dbpass, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Infof("Successfully connected!")
}
