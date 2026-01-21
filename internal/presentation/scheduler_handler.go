package presentation

import (
	"context"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/usecase"
)

type PointSchedulerHandler struct {
	pointReservationConfirm usecase.PointReservationConfirmUseCase
}

func NewPointSchedulerHandler(
	pointReservationConfirmUseCase usecase.PointReservationConfirmUseCase,
) *PointSchedulerHandler {
	return &PointSchedulerHandler{
		pointReservationConfirm: pointReservationConfirmUseCase,
	}
}

func (p *PointWorkerHandler) PointReserveConfirmScheduler(c context.Context, producer *rabbitmq.RabbitProducer) {

}
