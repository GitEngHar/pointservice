package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"pointservice/internal/domain"
	"pointservice/internal/usecase/tally"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	internalUri = "amqp://guest:guest@127.0.0.1:5672/"
)

// Rabbit MessageQueueの非同期通信
type RabbitProducer struct {
	producer *amqp.Connection
}

func NewRabbitProducer(producer *amqp.Connection) tally.Producer {
	return &RabbitProducer{
		producer: producer,
	}
}

func (r *RabbitProducer) PublishPoint(c context.Context, point domain.Point) error {
	if r.producer == nil {
		// TODO keploy検証ではRabbitMQの連携ができなかったので一時的にProducerのnilをkeploy用に許容する
		// TODO 暫定的に本番でnilならエラーにする
		return nil
	}
	ch, err := r.producer.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	queueName := "sampleQueue"
	queue, err := buildQueue(ch, queueName)
	if err != nil {
		return err
	}
	body, err := json.Marshal(point)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		c,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// ReservationMessage is the message format for reservation queue
type ReservationMessage struct {
	ReservationID  string `json:"reservation_id"`
	UserID         string `json:"user_id"`
	PointAmount    int    `json:"point_amount"`
	IdempotencyKey string `json:"idempotency_key"`
}

const reservationQueueName = "reservationQueue"

// 実際に予約メッセージをRabbitMQに「送信（Publish）」する関数。
func (r *RabbitProducer) PublishReservation(c context.Context, msg ReservationMessage) error {
	if r.producer == nil {
		return nil
	}
	ch, err := r.producer.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := buildQueue(ch, reservationQueueName) // 送り先となる「キュー（ポスト）」を用意する。
	if err != nil {
		return err
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// メッセージを「キュー」に投げる。
	return ch.PublishWithContext(
		c,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

const (
	retryInterval = time.Second * 5
	retryCount    = 3
)

func ConnectProducer(env string) (*amqp.Connection, func() error) {
	var retry int
	if env == "keploy" {
		return nil, nil
	}
	for {
		conn, err := amqp.Dial(internalUri)
		if err == nil {
			fmt.Println("connected to rabbitmq!!")
			return conn, conn.Close
		}
		fmt.Printf("failed to connect to rabbitmq!!! error %v. retry %v\n", err.Error(), retry+1)
		time.Sleep(retryInterval)
		if retry >= retryCount {
			panic(err)
		}
		retry++
	}
}
func buildQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	// TODO each mean
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
