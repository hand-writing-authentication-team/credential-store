package dao

import (
	"database/sql"
)

type TXManager struct {
	Tx *sql.Tx
}

func (txm *TXManager) Commit() (err error) {
	err = txm.Tx.Commit()
	return err
}

func (txm *TXManager) Rollback() (err error) {
	err = txm.Tx.Rollback()
	return err
}

func (txm *TXManager) End(err error) error {
	if err != nil {
		return txm.Commit()
	}
	return txm.Rollback()
}
