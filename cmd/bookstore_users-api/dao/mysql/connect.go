package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/conf"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/go-sql-driver/mysql"
)

func MakeConfig(cfg conf.Config) mysql.Config {
	var mysqlConf mysql.Config

	mysqlConf.Addr = fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port)
	mysqlConf.DBName = cfg.Database.Schema
	mysqlConf.User = cfg.Database.Uname
	mysqlConf.Passwd = ""
	return mysqlConf
}

func NewSqlClient(conf mysql.Config) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		conf.User, conf.Passwd, conf.Addr, conf.DBName)

	log.Println("connection :", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	mysql.SetLogger(logger.Logger{})
	log.Println("database successfully configured")
	return db
}
