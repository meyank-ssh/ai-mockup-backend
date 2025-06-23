package service

import (
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQ *amqp.Connection

func InitRabbitMQ() error {
	var err error
	RabbitMQ, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return err
	}
	
	return nil
}

func CloseRabbitMQ() {
	if RabbitMQ != nil {
		RabbitMQ.Close()
	}
} 