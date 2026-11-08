package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	topic       = "SC-D"
	poisonTopic = "SC-D.poison"
	maxRetries  = 3
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_USER / RABBITMQ_PASSWORD)
	if err != nil {
		panic(err)
	}
	cfg := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	publisher, err := amqp.NewPublisher(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	ctx := context.Background()

	// Subscribing to the poison topic declares its queue (so poison messages are not dropped)
	// and prints whatever lands there.
	poisoned, err := subscriber.Subscribe(ctx, poisonTopic)
	if err != nil {
		panic(err)
	}
	go func() {
		for msg := range poisoned {
			fmt.Printf("POISON QUEUE: id=%s payload=%s reason=%q\n",
				msg.UUID, string(msg.Payload), msg.Metadata.Get(middleware.ReasonForPoisonedKey))
			msg.Ack()
		}
	}()

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	poisonMW, err := middleware.PoisonQueue(publisher, poisonTopic)
	if err != nil {
		panic(err)
	}
	retryMW := middleware.Retry{
		MaxRetries:      maxRetries,
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2,
		MaxInterval:     2 * time.Second, // must be set: with 0 the backoff never grows
		OnRetryHook: func(n int, delay time.Duration) {
			fmt.Printf("retry #%d (next in %v)\n", n, delay)
		},
		Logger: logger,
	}
	// Order matters: PoisonQueue must be outermost so it sees the error after Retry gives up.
	router.AddMiddleware(poisonMW, retryMW.Middleware)

	router.AddConsumerHandler("sc-d-handler", topic, subscriber, func(msg *message.Message) error {
		if msg.Metadata.Get("poison") == "true" {
			fmt.Printf("attempt: id=%s FAILED\n", msg.UUID)
			return errors.New("cannot process message")
		}
		fmt.Printf("received: id=%s payload=%s\n", msg.UUID, string(msg.Payload))
		return nil
	})

	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
