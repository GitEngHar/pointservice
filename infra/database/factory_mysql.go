package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func ConnectDB() (*sql.DB, func() error) {
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
	if err = db.Ping(); err != nil {
		connectionRetry(err, mysqlConfig)
	}
	fmt.Println("db connected!!	")
	return db, func() error {
		return db.Close()
	}
}

func connectionRetry(err error, mysqlConfig config) {
	var connectRetryNum int
	fmt.Println("db connection failed retry start")
	for err != nil {
		time.Sleep(mysqlConfig.retryInterval)
		if connectRetryNum > mysqlConfig.maxRetryConnection {
			log.Fatalf("db connection time out (3min) %s", err)
		}
		connectRetryNum++
		fmt.Printf("db connection retry...%d \n", connectRetryNum)
	}
}
