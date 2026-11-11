package main

import (
	"context"
	"fmt"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	rabbit "github.com/rabbitmq/amqp091-go"
)

const (
	topic = "SC-D.quorum"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_USER / RABBITMQ_PASSWORD)
	if err != nil {
		panic(err)
	}
	cfg := amqp.NewDurableQueueConfig(amqpURI)
	cfg.Marshaler = amqp.DefaultMarshaler{
		PreprocessDelivery: func(d rabbit.Delivery) rabbit.Delivery {
			for k, v := range d.Headers {
				if _, ok := v.(string); !ok {
					d.Headers[k] = fmt.Sprint(v)
				}
			}
			return d
		},
	}
	cfg.Queue.Arguments = rabbit.Table{
		"x-queue-type":              "quorum",
		"x-delivery-limit":          3,
		"x-dead-letter-exchange":    "sc-d.dlx",
		"x-dead-letter-routing-key": "sc-d.dead",
	}

	subscriber, err := amqp.NewSubscriber(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	ctx := context.Background()

	messages, err := subscriber.Subscribe(ctx, topic)
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		if msg.Metadata.Get("poison") == "true" {
			fmt.Printf("id=%s FAILED\n", msg.UUID)
			//os.Exit(1)
			msg.Nack()
			continue
		}
		fmt.Printf("received: id=%s payload=%s\n", msg.UUID, string(msg.Payload))
		msg.Ack()
	}
}
