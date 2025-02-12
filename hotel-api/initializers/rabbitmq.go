package initializers

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection
var RabbitMQChannel *amqp.Channel

// Conectar a RabbitMQ
func ConnectRabbitMQ() error {
	var err error
	RabbitMQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return err
	}

	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return err
	}

	// Asegurarse de que la conexión y el canal están activos
	if RabbitMQConn.IsClosed() {
		log.Fatal("RabbitMQ connection is closed")
		return fmt.Errorf("RabbitMQ connection is closed")
	}

	if RabbitMQChannel == nil {
		log.Fatal("RabbitMQ channel is nil")
		return fmt.Errorf("RabbitMQ channel is nil")
	}

	log.Println("Successfully connected to RabbitMQ")
	return nil
}
