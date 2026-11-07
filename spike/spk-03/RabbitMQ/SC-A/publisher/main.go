package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const topic = "hello"

func main() {
	logger := watermill.NewStdLogger(false, false)

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_USER / RABBITMQ_PASSWORD)
	if err != nil {
		panic(err)
	}

	publisher, err := amqp.NewPublisher(amqp.NewDurableQueueConfig(amqpURI), logger)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	for i := 1; i <= 100; i++ {
		uuid := watermill.NewUUID()
		payload := fmt.Sprintf("message #%d", i)
		msg := message.NewMessage(uuid, []byte(payload))

		if err := publisher.Publish(topic, msg); err != nil {
			panic(err)
		}
		fmt.Printf("published: id=%s, payload=%s\n", uuid, payload)
	}
}
