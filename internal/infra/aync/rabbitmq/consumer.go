package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	externalUri = "amqp://guest:guest@rabbitmq:5672/"
)

// Rabbit MessageQueueの非同期通信
type RabbitConsumer struct {
}

func NewRabbitConsumer() *RabbitConsumer {
	return &RabbitConsumer{}
}

func (r *RabbitConsumer) GetSumPoint(_ context.Context) error {
	conn, err := amqp.Dial(externalUri)
	if err != nil {
		return err
	}
	defer conn.Close()
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
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("waiting for messages...")

	// 5. 受信
	for msg := range msgs {
		log.Printf("received: %s", msg.Body)
		// 処理成功とみなしてAckをする
		err = msg.Ack(false)
		if err != nil {
			return err
		}
	}
	return nil
}
