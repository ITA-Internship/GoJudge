package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	topic    = "SC-C"
	messages = 50 // unique keys; each is published twice (original + producer duplicate)
)

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

	// The idempotency key identifies the logical operation (one submission)
	// and travels with the message in its metadata.
	keys := make([]string, messages)
	for i := range keys {
		keys[i] = watermill.NewUUID()
	}

	publish := func(kind string, i int) {
		payload := fmt.Sprintf("submission #%d", i+1)
		msg := message.NewMessage(watermill.NewUUID(), []byte(payload)) // new message UUID every time
		msg.Metadata.Set("idempotency_key", keys[i])

		if err := publisher.Publish(topic, msg); err != nil {
			panic(err)
		}
		fmt.Printf("published(%s): n=%d key=%s id=%s\n", kind, i+1, keys[i], msg.UUID)
	}

	// Round 1: the original messages.
	for i := range keys {
		publish("original", i)
	}
	// Round 2: producer-side duplicates. Same key, different message UUID
	// (like a retried request that publishes again). The broker cannot tell
	// these apart from new messages, only the idempotency key can.
	for i := range keys {
		publish("producer-duplicate", i)
	}

	fmt.Printf("SUMMARY published=%d unique_keys=%d\n", 2*messages, messages)
}
