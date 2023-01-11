package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB
var SQLDB *sql.DB

func Connect(dsn string) {
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 获取底层的 sqlDB
	SQLDB = DB.DB
	if err != nil {
		fmt.Println(err.Error())
	}
}
