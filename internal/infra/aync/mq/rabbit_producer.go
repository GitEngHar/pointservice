package mq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"pointservice/internal/domain"
	"pointservice/internal/usecase/tally"
)

const (
	internalUri = "amqp://guest:guest@rabbitmq:5672/"
)

// Rabbit MessageQueueの非同期通信
type RabbitProducer struct {
}

func NewRabbitProducer() tally.Producer {
	return &RabbitProducer{}
}

func (r *RabbitProducer) PublishPoint(c context.Context, point domain.Point) error {
	conn, err := amqp.Dial(internalUri)
	if err != nil {
		return err
	}
	defer conn.Close()
	//TODO channel?
	ch, err := conn.Channel()
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
