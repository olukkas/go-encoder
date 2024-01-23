package queue

import (
	"fmt"
	"github.com/olukkas/go-encoder/framework/utils"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitMQ struct {
	User              string
	Password          string
	Host              string
	Port              string
	Vhost             string
	ConsumerQueueName string
	ConsumerName      string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	rabbitMqArgs := amqp.Table{
		"x-dead-letter-exchange": os.Getenv("RABBITMQ_DLX"),
	}

	return &RabbitMQ{
		User:              os.Getenv("RABBITMQ_USER"),
		Password:          os.Getenv("RABBITMQ_PASS"),
		Host:              os.Getenv("RABBITMQ_HOST"),
		Port:              os.Getenv("RABBITMQ_PORT"),
		Vhost:             os.Getenv("RABBITMQ_VHOST"),
		ConsumerQueueName: os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		ConsumerName:      os.Getenv("RABBITMQ_CONSUMER_NAME"),
		AutoAck:           false,
		Args:              rabbitMqArgs,
	}
}

func (r *RabbitMQ) Connect() *amqp.Channel {
	dns := fmt.Sprintf("amqp://%s:%s@%s:%s%s", r.User, r.Password, r.Host, r.Port, r.Vhost)

	conn, err := amqp.Dial(dns)
	utils.FailOnError(err, "failed to connect to RabbitMQ")

	r.Channel, err = conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	return r.Channel
}

func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,
		true,
		false,
		false,
		false,
		r.Args,
	)
	utils.FailOnError(err, "failed to declare a queue")

	incomingMessages, err := r.Channel.Consume(
		q.Name,
		r.ConsumerName,
		r.AutoAck,
		false,
		false,
		false,
		r.Args,
	)
	utils.FailOnError(err, "Failed to register a consumer")

	go func() {
		for message := range incomingMessages {
			log.Println("Incoming new Message")
			messageChannel <- message
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

func (r *RabbitMQ) Notify(message, contentType, exchange, routingKey string) error {
	return r.Channel.Publish(
		exchange, routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)
}
