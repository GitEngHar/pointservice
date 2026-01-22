package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

//const (
//	reservationQueueName = "reservationQueue"
//)

func NewQueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		reservationQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func NewConsume(ch *amqp.Channel, queue amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		queue.Name,
		"",    // consumer tag
		false, // auto-ack (manual ack for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
}
