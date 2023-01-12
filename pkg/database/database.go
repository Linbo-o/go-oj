package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
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

func IsExist(table, filed, value string) bool {
	sql_ := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?", filed, table, filed)
	log.Println(sql_)
	var val string
	err := SQLDB.QueryRow(sql_, value).Scan(&val)
	if err != nil {
		return false
	}
	return true
}
