package mq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "tp.tasks.exchange"
)

//Message defines a message request object
type Message struct {
	CorrelationID  string
	Request        interface{}
	RequestAction  string
	Response       interface{}
	ResponseWorker string
	ResponseCode   int
}

//MessageQueue defines a message queue object
type MessageQueue struct {
	channel *amqp.Channel
}

//NewMQ creates a new message queue object
func NewMQ() *MessageQueue {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	return &MessageQueue{
		channel: ch,
	}
}

//NewQueue creates a new message queue and bind to the exchange
func (mq *MessageQueue) NewQueue(name string, routingKey string) {
	q, err := mq.channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to create a queue: "+name)

	err = mq.channel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
}

//PublishMessage creates a new task to a queue
func (mq *MessageQueue) PublishMessage(message Message, routingKey string, corrID string, replyTo string) {
	json, err := json.Marshal(message)
	failOnError(err, "Failed to parse message to JSON")

	err = mq.channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       replyTo,
			Body:          json,
		})
	failOnError(err, "Failed to publish a message")
}

//GetMessages sent to the queue
func (mq *MessageQueue) GetMessages(queueName string) <-chan amqp.Delivery {
	msgs, err := mq.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
