package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        dat,
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
}

func PublishGob[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(val)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
}
