package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	topic          = "SC-D"
	poisonMessages = 3
)

// Start the consumer first: the queue is declared by the subscriber, not the publisher.
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

	for range poisonMessages {
		// 1 poison message between 2 normal ones.
		for i, poison := range []bool{false, true, false} {
			uuid := watermill.NewUUID()
			payload := fmt.Sprintf("message #%d", i+1)
			msg := message.NewMessage(uuid, []byte(payload))
			if poison {
				msg.Metadata.Set("poison", "true")
			}

			if err := publisher.Publish(topic, msg); err != nil {
				panic(err)
			}
			fmt.Printf("published: id=%s, payload=%s, poison=%t\n", uuid, payload, poison)
		}
	}
}
