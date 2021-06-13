package users_db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/config"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/logger"
	"github.com/go-sql-driver/mysql"
)

func MakeMySQLConfig(cfg config.Config) mysql.Config {
	var mysqlConf mysql.Config

	mysqlConf.Addr = fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port)
	mysqlConf.DBName = cfg.Database.Schema
	mysqlConf.User = cfg.Database.Uname
	mysqlConf.Passwd = ""
	return mysqlConf
}

func ProvideSqlClient(dbconf mysql.Config) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		dbconf.User, dbconf.Passwd, dbconf.Addr, dbconf.DBName)

	log.Println("connection :", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	mysql.SetLogger(logger.Logger{})
	log.Println("database successfully configured")
	return db
}
