package infra

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

type dbInfo struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
}

func ConnectDB() (*sql.DB, func() error) {
	connectDBInfo := dbInfo{
		dbUser:     os.Getenv("dbUser"),
		dbPassword: os.Getenv("dbPassword"),
		dbHost:     os.Getenv("dbHost"),
		dbPort:     os.Getenv("dbPort"),
		dbName:     os.Getenv("dbName"),
	}
	var connectRetryNum int
	connectDBUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", connectDBInfo.dbUser, connectDBInfo.dbPassword, connectDBInfo.dbHost, connectDBInfo.dbPort, connectDBInfo.dbName)
	db, err := sql.Open("mysql", connectDBUri)
	if err != nil {
		panic(err)
	}

	// ConnectionPool Setting
	// TODO: to env params
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	err = db.Ping()
	for err != nil {
		time.Sleep(3 * time.Second)
		if connectRetryNum > 20 {
			log.Fatalf("db connection time out (3min) %s", err)
		}
		connectRetryNum++
	}
	fmt.Println("db connected!!")
	return db, func() error {
		// db close closure
		return db.Close()
	}
}
