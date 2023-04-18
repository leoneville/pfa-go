package main

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

func GenerateOrders() Order {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	return Order{
		ID:    uuid.New().String(),
		Price: seed.Float64() * 100,
	}
}

func Notify(ch *amqp.Channel, order Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"amq.direct", // exchange
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	for i := 0; i < 100000000; i++ {
		order := GenerateOrders()
		if err := Notify(ch, order); err != nil {
			panic(err)
		}
		//fmt.Println(order)
	}
}
