package rabbitmq

import (
	"context"
	"encoding/json"
	"pointservice/internal/infra/aync/dto"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	internalUri = "amqp://guest:guest@rabbitmq:5672/"
)

// RabbitProducer Rabbit MessageQueueの非同期通信
type RabbitProducer struct {
	producer *amqp.Connection
}

func NewRabbitProducer(producer *amqp.Connection) *RabbitProducer {
	return &RabbitProducer{
		producer: producer,
	}
}

const reservationQueueName = "reservationQueue"

// 実際に予約メッセージをRabbitMQに「送信（Publish）」する関数。
func (r *RabbitProducer) PublishReservation(c context.Context, msg dto.ReservationMessage) error {
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
