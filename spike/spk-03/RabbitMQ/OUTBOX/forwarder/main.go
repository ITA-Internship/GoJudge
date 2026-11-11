// SC-E forwarder: the relay of the transactional outbox (Watermill's forwarder component).
//
// It reads the outbox table (SQL subscriber on forwarderTopic), unwraps every envelope
// and publishes the original message to the destination topic in RabbitMQ. A row is
// marked as processed only after the publish succeeded (Ack), so a crash or an error
// before that point means the row is sent again: at-least-once delivery.
//
// -crash-after-publish kills the process right after the first message is published
// and confirmed by the broker, but before it is acknowledged. It imitates the window
// that every relay has between "sent" and "marked as sent".
//
// Run order: consumer first (it declares the queue), then this program.
// Stop with Ctrl+C.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/OUTBOX/outlog"
	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const forwarderTopic = "outbox_forwarder" // outbox table: "watermill_outbox_forwarder"

var out io.Writer = os.Stdout

func main() {
	crashAfterPublish := flag.Bool("crash-after-publish", false, "kill the process after the first message is published, before the Ack")
	logPath := flag.String("log", "", "also write the output to this log file")
	flag.Parse()

	var closeLog func()
	out, closeLog = outlog.Open(*logPath)
	defer closeLog()
	logger := watermill.NewStdLoggerWithOut(out, false, false)

	dsn, err := rmqenv.PostgresDSN() // credentials come from .env (POSTGRES_*)
	must(err)
	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_*)
	must(err)

	// In: the outbox table. InitializeSchema creates the outbox and the offsets tables if they are missing.
	sqlSubscriber, err := watermillSQL.NewSubscriber(db, watermillSQL.SubscriberConfig{
		SchemaAdapter:    watermillSQL.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)
	must(err)

	// Out: RabbitMQ. ConfirmDelivery makes Publish wait for the broker to confirm the message,
	// without it a message could be lost after the forwarder has already acknowledged the row.
	cfg := amqp.NewDurableQueueConfig(amqpURI)
	cfg.Publish.ConfirmDelivery = true
	amqpPublisher, err := amqp.NewPublisher(cfg, logger)
	must(err)
	defer amqpPublisher.Close()

	fwd, err := forwarder.NewForwarder(
		sqlSubscriber,
		loggedPublisher{Publisher: amqpPublisher, crashAfterPublish: *crashAfterPublish},
		logger,
		forwarder.Config{ForwarderTopic: forwarderTopic},
	)
	must(err)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Fprintf(out, "forwarder started: outbox table -> RabbitMQ (crash-after-publish=%v), Ctrl+C to stop\n", *crashAfterPublish)
	must(fwd.Run(ctx)) // blocks until Ctrl+C
	fmt.Fprintln(out, "forwarder stopped")
}

// loggedPublisher prints every message that reached RabbitMQ and, if asked,
// kills the process before the forwarder can acknowledge the outbox row.
type loggedPublisher struct {
	message.Publisher
	crashAfterPublish bool
}

func (p loggedPublisher) Publish(topic string, msgs ...*message.Message) error {
	if err := p.Publisher.Publish(topic, msgs...); err != nil {
		fmt.Fprintf(out, "publish FAILED, the row stays in the outbox and is retried: %v\n", err)
		return err
	}
	for _, m := range msgs {
		fmt.Fprintf(out, "forwarded: submission_id=%s message_uuid=%s -> queue %s\n", m.Metadata.Get("submission_id"), m.UUID, topic)
	}
	if p.crashAfterPublish {
		fmt.Fprintln(out, "CRASH: killed after the publish, before the Ack of the outbox row")
		os.Exit(1)
	}
	return nil
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
		panic(err)
	}
}
