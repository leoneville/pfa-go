package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leoneville/pfa-go/internal/order/infra/database"
	"github.com/leoneville/pfa-go/internal/order/usecase"
	"github.com/leoneville/pfa-go/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	maxWaiters := 3
	wg := sync.WaitGroup{}

	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	uc := usecase.NewCalculateFinalPriceUseCase(repository)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	out := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, out)

	wg.Add(maxWaiters)
	for i := 0; i < maxWaiters; i++ {
		defer wg.Done()
		go worker(out, uc, i)
	}
	wg.Wait()

}

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerId int) {
	for msg := range deliveryMessage {
		var input usecase.OrderInputDTO
		if err := json.Unmarshal(msg.Body, &input); err != nil {
			fmt.Println("Error unmarshalling message", err)
		}
		input.Tax = 10.0
		_, err := uc.Execute(input)
		if err != nil {
			fmt.Println("Error unmarshalling message", err)
		}
		msg.Ack(false)
		fmt.Println("Worker", workerId, "Processed order", input.ID)
	}
}
