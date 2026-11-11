package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	rabbit "github.com/rabbitmq/amqp091-go"
)

const (
	topic      = "SC-D.2"
	maxRetries = 3
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_USER / RABBITMQ_PASSWORD)
	if err != nil {
		panic(err)
	}
	cfg := amqp.NewDurableQueueConfig(amqpURI)
	cfg.Consume.NoRequeueOnNack = true
	cfg.Queue.Arguments = rabbit.Table{
		"x-dead-letter-exchange":    "sc-d.dlx",
		"x-dead-letter-routing-key": "sc-d.dead",
	}

	// Subscriber for the main queue: nack without requeue -> broker dead-letters to sc-d.dlx.
	subscriber, err := amqp.NewSubscriber(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	ctx := context.Background()

	router, err := message.NewRouter(message.RouterConfig{}, logger)
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
	router.AddMiddleware(retryMW.Middleware)

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
