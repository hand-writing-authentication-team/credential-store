package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PGDBInstance struct {
	conn *sql.DB
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
