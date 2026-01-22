package presentation

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"pointservice/internal/usecase"
)

type PointWorkerHandler struct {
	pointReservationAdd usecase.PointReservationAddUseCase
}

func NewPointWorkerHandler(
	pointReservationAddUseCase usecase.PointReservationAddUseCase,
) *PointWorkerHandler {
	return &PointWorkerHandler{
		pointReservationAdd: pointReservationAddUseCase,
	}
}

func (p *PointWorkerHandler) PointReserveWorker(c context.Context, msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		if er := p.pointReservationAdd.Execute(c, msg); er != nil {
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
