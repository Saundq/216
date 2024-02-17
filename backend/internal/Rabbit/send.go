package Rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var err error

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Connect(RabbitMqString string) {
	Conn, err = amqp.Dial(RabbitMqString)
	log.Println("RabbitMQ connection init")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer Conn.Close()

}

func Read() *amqp.Channel {
	channelRabbitMQ, err := Conn.Channel()

	if err != nil {
		log.Println(err)
		panic(err)
	}
	//defer channelRabbitMQ.Close()
	return channelRabbitMQ
}

func Send(channel string, message string) {
	ch, err := Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel, //name
		false,   //durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body := message
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	//log.Printf(" [x] Sent %s\n", body)
}
