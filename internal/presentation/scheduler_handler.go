package presentation

import (
	"context"
	"log"
	"pointservice/internal/infra/aync/rabbitmq"
	"pointservice/internal/usecase"
	"time"
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

const scanInterval = 10 * time.Second

// PointReserveConfirmScheduler ドメイン層（ビジネスロジック）とUI層（インターフェース）の橋渡しをする役割
func (p *PointSchedulerHandler) PointReserveConfirmScheduler(c context.Context, producer *rabbitmq.RabbitProducer) {
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	log.Printf("Scheduler running, scanning every %v\n", scanInterval)

	for {
		select {
		case <-ticker.C: // 10秒経ったら、中の処理をやる。
			if err := p.pointReservationConfirm.Execute(c, producer); err != nil {
				log.Printf("Error during scan: %v\n", err)
			}
		case <-c.Done(): // 終了合図が来たら、終了する。
			log.Printf("shutting down...\n")
			return
		}
	}
}
