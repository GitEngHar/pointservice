package main

import (
	"context"
	"os"
	"pointservice/internal/infra/aync/mq"
	"pointservice/internal/usecase/tally"
)

// batchの実行
func main() {
	consumer := mq.NewRabbitConsumer()
	useCase := tally.NewPointTally(consumer)
	err := useCase.Execute(context.Background())
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
