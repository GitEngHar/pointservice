package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func ConnectDB() (*sql.DB, func() error) {
	var connectRetryNum int
	mysqlConfig := newConfigMysql()
	connectDBUri := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		mysqlConfig.user, mysqlConfig.password, mysqlConfig.host, mysqlConfig.port, mysqlConfig.database,
	)
	db, err := sql.Open("mysql", connectDBUri)
	if err != nil {
		log.Fatalf("sql open failed: %s", err)
	}
	db.SetMaxOpenConns(mysqlConfig.maxOpenConnections)
	db.SetMaxIdleConns(mysqlConfig.maxIdleConnections)
	db.SetConnMaxLifetime(mysqlConfig.maxOpenConnectionTime)
	for {
		time.Sleep(mysqlConfig.retryInterval)
		err := db.Ping()
		if err == nil {
			break
		}
		connectRetryNum++
		if connectRetryNum > mysqlConfig.maxRetryConnection {
			log.Fatalf("db connection time out (3min): %s", err)
		}
	}
	fmt.Println("db connected!!	")
	return db, func() error {
		return db.Close()
	}
}
