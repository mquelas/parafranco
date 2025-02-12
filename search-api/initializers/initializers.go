package initializers

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var RabbitMQChannel *amqp091.Channel

// Inicializar RabbitMQ
func InitRabbitMQ() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	RabbitMQChannel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	fmt.Println("RabbitMQ connected")
}
