package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable   SimpleQueueType = 0
	SimpleQueueTransient SimpleQueueType = 1
)

type Acktype int

const (
	Ack         Acktype = 0
	NackReque   Acktype = 1
	NackDiscard         = 2
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		queueName,                             // name
		simpleQueueType == SimpleQueueDurable, // durable
		simpleQueueType != SimpleQueueDurable, // delete when unused
		simpleQueueType != SimpleQueueDurable, // exclusive
		false,                                 // no-wait
		nil,                                   // argument
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	go func() {
		defer ch.Close()

		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("Could not unmarshal message: %v\n", err)
				continue
			}

			acktype := handler(target)
			switch acktype {
			case Ack:
				log.Println("Ack")
				msg.Ack(false)
			case NackReque:
				log.Println("Nack Reque")
				msg.Nack(false, true)
			case NackDiscard:
				log.Println("Nack Discard")
				msg.Nack(false, false)
			}
		}
	}()

	return nil
}
