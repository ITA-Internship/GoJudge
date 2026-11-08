package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	topic      = "SC-C"
	deliveries = 150
)

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

	// Final state: idempotency key -> how many times the side effect really ran.
	// Messages arrive one at a time on this channel, so no mutex is needed.
	applied := map[string]int{}
	received, skipped := 0, 0

	for msg := range messages {
		received++
		key := msg.Metadata.Get("idempotency_key")

		if applied[key] > 0 {
			// Second delivery of an already processed message: do nothing, Ack for real.
			fmt.Printf("skipped duplicate: key=%s id=%s\n", key, msg.UUID)
			skipped++
			msg.Ack()
		} else {
			// First delivery: do the side effect, then "lose the Ack" with a Nack.
			// The broker puts the same message back and delivers it again.
			applied[key]++
			fmt.Printf("applied: key=%s id=%s\n", key, msg.UUID)
			msg.Nack()
		}

		if received >= deliveries {
			break
		}
	}

	appliedTotal, maxPerKey := 0, 0
	for _, n := range applied {
		appliedTotal += n
		if n > maxPerKey {
			maxPerKey = n
		}
	}

	fmt.Printf("SUMMARY deliveries=%d keys=%d applied_total=%d skipped=%d max_applied_per_key=%d\n",
		received, len(applied), appliedTotal, skipped, maxPerKey)
}
