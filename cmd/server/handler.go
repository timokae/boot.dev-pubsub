package main

import (
	"fmt"
	"log"

	"github.com/timokae/learn-pub-sub-starter/internal/gamelogic"
	"github.com/timokae/learn-pub-sub-starter/internal/pubsub"
	"github.com/timokae/learn-pub-sub-starter/internal/routing"
)

func handlerLog() func(gameLog routing.GameLog) pubsub.Acktype {
	return func(gameLog routing.GameLog) pubsub.Acktype {
		defer fmt.Println("> ")

		err := gamelogic.WriteLog(gameLog)
		if err != nil {
			log.Printf("error writing log: %v\n", err)
		}
		return pubsub.Ack
	}
}
