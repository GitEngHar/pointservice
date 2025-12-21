package mysql

import (
	"log"
	"os"
	"strconv"
	"time"
)

type config struct {
	host                  string
	database              string
	port                  string
	user                  string
	password              string
	maxOpenConnections    int
	maxIdleConnections    int
	maxRetryConnection    int
	maxOpenConnectionTime time.Duration
	retryInterval         time.Duration
}

func newConfigMysql() config {
	maxOpenConnections, err := strconv.Atoi(os.Getenv("POINT_MYSQL_MAX_OPEN_CONNECTIONS"))
	maxIdleConnections, err := strconv.Atoi(os.Getenv("POINT_MYSQL_MAX_IDLE_CONNECTIONS"))

	defer func() {
		if err != nil {
			log.Fatalf("Invalid environment variable value: %s", err)
		}
	}()

	return config{
		host:                  os.Getenv("POINT_MYSQL_HOST"),
		database:              os.Getenv("POINT_MYSQL_DATABASE"),
		port:                  os.Getenv("POINT_MYSQL_PORT"),
		user:                  os.Getenv("POINT_MYSQL_USER"),
		password:              os.Getenv("POINT_MYSQL_PASSWORD"),
		maxOpenConnections:    maxOpenConnections,
		maxIdleConnections:    maxIdleConnections,
		maxRetryConnection:    5,
		maxOpenConnectionTime: 15 * time.Minute,
		retryInterval:         5 * time.Second,
	}
}
