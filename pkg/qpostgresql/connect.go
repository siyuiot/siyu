package qpostgresql

import (
	"database/sql"
	"fmt"

	"github.com/go-xorm/xorm"
)

func InitPg(url string, maxopen, maxidle int) (db *sql.DB, xdb *xorm.Engine, err error) {
	db, err = sql.Open("postgres", url)
	if err != nil {
		fmt.Printf("conn db faild,err=%v\n", err)
		return
	}
	db.SetMaxOpenConns(maxopen / 2)
	db.SetMaxIdleConns(maxidle / 2)
	err = db.Ping()
	if err != nil {
		fmt.Printf("db ping faild,err=%v\n", err)
		return
	}

	return
}
