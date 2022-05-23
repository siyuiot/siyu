package qpostgresql

import (
	"database/sql"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
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
func InitGorm(url string, maxopen, maxidle int) (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", url)
	if err != nil {
		return
	}
	sqlDb := db.DB()
	sqlDb.SetMaxOpenConns(maxopen / 2)
	sqlDb.SetMaxIdleConns(maxidle / 2)
	err = sqlDb.Ping()
	if err != nil {
		fmt.Printf("xdb ping faild,err=%v\n", err)
		return
	}
	db.LogMode(true)
	return
}
