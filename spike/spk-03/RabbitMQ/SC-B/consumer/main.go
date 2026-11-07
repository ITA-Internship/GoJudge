package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const topic = "SC-B"

func main() {
	logger := watermill.NewStdLogger(false, false)

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_USER / RABBITMQ_PASSWORD)
	if err != nil {
		panic(err)
	}

	subscriber, err := amqp.NewSubscriber(amqp.NewDurableQueueConfig(amqpURI), logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	messages, err := subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		fmt.Printf("received: id=%s payload=%s\n", msg.UUID, string(msg.Payload))
		os.Exit(1)
		msg.Ack()
		subscriber.Close()
	}
}
