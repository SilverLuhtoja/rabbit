package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbit_con_path := "amqp://guest:guest@localhost:5672/"
	signalChan := make(chan os.Signal, 1) // signal channel to listen connections

	connection, err := amqp.Dial(rabbit_con_path)
	if err != nil {
		log.Fatal("Coulnt connect to rabbit server")
	}
	defer connection.Close()
	fmt.Println("Server: Connection was successful...")

	channel, err := connection.Channel()
	if err != nil {
		log.Fatal("Couldn't start channel")
	}

	pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	gamelogic.PrintClientHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		firstWord := words[0]

		switch firstWord {
		case "pause":
			fmt.Println("Sending pause message")
			pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		}

	}

	// liste ctrl + c
	signal.Notify(signalChan, os.Interrupt)
	signal_message := <-signalChan
	fmt.Printf("Signal received: %s\n Programm is shutting down... \n", signal_message)
}
