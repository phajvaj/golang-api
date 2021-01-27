package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func SetConnection() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbHost := viper.GetString("dbms.host")
	dbPort := viper.GetString("dbms.port")
	dbName := viper.GetString("dbms.dbname")
	dbUser := viper.GetString("dbms.username")
	dbPass := viper.GetString("dbms.password")

	db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8", dbUser, dbPass, dbHost, dbPort, dbName))
	//defer db.Close()
	if err != nil {
		fmt.Println("connect fail")
	} else {
		fmt.Println("connect success")
	}
	return
}
