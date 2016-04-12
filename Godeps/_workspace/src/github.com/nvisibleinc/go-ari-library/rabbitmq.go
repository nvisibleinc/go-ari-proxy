package ari

import (
	"github.com/streadway/amqp"
)

type rabbitmqConfig struct {
	URL string `json:"url"`
}
type RabbitMQ struct {
	config       rabbitmqConfig
	producerConn *amqp.Connection
	consumerConn *amqp.Connection
}

func (r *RabbitMQ) InitBus(config interface{}) error {
	var err error
	c := config.(map[string]interface{})
	for key, value := range c {
		switch key {
		case "url":
			r.config.URL = value.(string)
		}
	}

	r.producerConn, err = amqp.Dial(r.config.URL)
	if err != nil {
		return err
	}
	r.consumerConn, err = amqp.Dial(r.config.URL)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) StartProducer(topic string) (chan []byte, error) {
	c := make(chan []byte)
	channel, err := r.producerConn.Channel()
	_, err = channel.QueueDeclare(
		topic, // name of queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		true,  // nowait
		nil)   // arguments

	go func(channel *amqp.Channel, messages chan []byte) {
		for message := range messages {
			channel.Publish(
				"", // exchange, for now always using the default exchange
				topic,
				false,
				false,
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "application/json",
					ContentEncoding: "",
					Body:            message,
					DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
					Priority:        0,              // 0-9
				})
		}
	}(channel, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *RabbitMQ) StartConsumer(topic string) (chan []byte, error) {
	c := make(chan []byte)
	channel, err := r.consumerConn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		topic, // name of queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		true,  // nowait
		nil)   // arguments

	if err != nil {
		return nil, err
	}
	deliveries, err := channel.Consume(queue.Name, "", false, false, true, true, nil)
	if err != nil {
		return nil, err
	}
	go func(deliveries <-chan amqp.Delivery, c chan []byte) {
		for d := range deliveries {
			c <- d.Body
			d.Ack(false) // false does *not* mean don't acknowledge, see library docs for details
		}
	}(deliveries, c)

	return c, nil
}

func (r *RabbitMQ) TopicExists(topic string) bool {
	return true
}
