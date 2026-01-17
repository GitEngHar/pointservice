package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const keployEnvironment = "keploy"

type Connection struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewConnection(isChannel bool, env string) *Connection {
	conn := connect(env)
	if !isChannel {
		return &Connection{
			Conn: conn,
			Ch:   nil,
		}
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &Connection{
		Conn: conn,
		Ch:   ch,
	}
}

func connect(env string) *amqp.Connection {
	if env == keployEnvironment {
		return nil
	}
	for range retryCount {
		conn, err := amqp.Dial(internalUri)
		if err == nil {
			fmt.Println("connected to rabbitmq!!")
			return conn
		}
		fmt.Printf("failed to connect to rabbitmq: %v", err)
		time.Sleep(retryInterval)
	}
	panic("failed to connect to rabbitmq")
}
