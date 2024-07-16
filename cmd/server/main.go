package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/timokae/learn-pub-sub-starter/internal/gamelogic"
	"github.com/timokae/learn-pub-sub-starter/internal/pubsub"
	"github.com/timokae/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")

	rabbitConnString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	log.Println("Successfully connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, gameLogQueue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Queue %v declared and bound!\n", gameLogQueue.Name)

	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		gameLogQueue.Name,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLog(),
	)
	if err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Sending pause message")
			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
		case "resume":
			fmt.Println("Sending resume message")
			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			fmt.Println("Quitting...")
			return
		default:
			fmt.Println("Unknown command. Try again.")
		}
	}
}
