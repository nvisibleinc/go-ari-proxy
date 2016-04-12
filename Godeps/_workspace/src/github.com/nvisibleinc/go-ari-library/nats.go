package ari

import (
	"github.com/apcera/nats"
)

type natsConfig struct {
	URL   string `json:"url"`
	Queue string `json:"queue"`
}
type NATS struct {
	config     natsConfig
	connection *nats.Conn
	encoder    *nats.EncodedConn
}

func (n *NATS) InitBus(config interface{}) error {
	var err error
	c := config.(map[string]interface{})
	for key, value := range c {
		switch key {
		case "url":
			n.config.URL = value.(string)
		case "queue":
			n.config.Queue = value.(string)
		}
	}

	n.connection, err = nats.Connect(n.config.URL)
	if err != nil {
		return err
	}
	n.encoder, err = nats.NewEncodedConn(n.connection, "default")
	if err != nil {
		return err
	}
	return nil
}

func (n *NATS) StartProducer(topic string) (chan []byte, error) {
	c := make(chan []byte)
	err := n.encoder.BindSendChan(topic, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (n *NATS) StartConsumer(topic string) (chan []byte, error) {
	c := make(chan []byte)
	_, err := n.encoder.BindRecvQueueChan(topic, n.config.Queue, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (n *NATS) TopicExists(topic string) bool {
	return true
}
