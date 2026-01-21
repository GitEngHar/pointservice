package presentation

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"pointservice/internal/usecase"
)

type PointWorkerHandler struct {
	addReservationPointUseCase usecase.AddReservationPointUseCase
}

func NewPointWorkerHandler(
	addReservationPointUseCase usecase.AddReservationPointUseCase,
) *PointWorkerHandler {
	return &PointWorkerHandler{
		addReservationPointUseCase,
	}
}

func (p *PointWorkerHandler) PointReserveWorker(c context.Context, msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		if er := p.addReservationPointUseCase.Execute(c, msg); er != nil {
			log.Printf("Error processing message: %v", er)
			// このメッセージの処理に失敗しました」とキュー（例：RabbitMQ など）に通知する
			if nackErr := msg.Nack(false, true); nackErr != nil {
				log.Printf("Failed to nack message: %v", nackErr)
			}
			continue
		}
		// Ack on success
		if ackErr := msg.Ack(false); ackErr != nil {
			log.Printf("Failed to ack message: %v", ackErr)
		}
	}
}
