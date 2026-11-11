// SC-E consumer / checker: receives the events from RabbitMQ and, when it stops, compares them
// with the business rows that really are in PostgreSQL.
//
//	data_without_message  rows in spk03_submissions that never produced a message
//	message_without_data  messages whose submission is not in the database
//
// Both must be 0 for the outbox to be consistent. duplicates counts repeated deliveries
// of the same submission (allowed: the outbox is at-least-once).
//
// It stops on Ctrl+C, or by itself after -idle without a message (used by runall).
// Start it first, it declares the queue.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	_ "github.com/lib/pq"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/OUTBOX/outlog"
	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)
 
const targetTopic = "outbox_submissions"
 
var out io.Writer = os.Stdout
 
func main() {
	logPath := flag.String("log", "", "also write the output to this log file")
	idle := flag.Duration("idle", 0, "stop after this long without a message (0 = only on Ctrl+C)")
	flag.Parse()
 
	var closeLog func()
	out, closeLog = outlog.Open(*logPath)
	defer closeLog()
	logger := watermill.NewStdLoggerWithOut(out, false, false)
 
	amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_*)
	must(err)
	dsn, err := rmqenv.PostgresDSN() // credentials come from .env (POSTGRES_*)
	must(err)
	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
 
	subscriber, err := amqp.NewSubscriber(amqp.NewDurableQueueConfig(amqpURI), logger)
	must(err)
	defer subscriber.Close()
 
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
 
	messages, err := subscriber.Subscribe(ctx, targetTopic)
	must(err)
	fmt.Fprintf(out, "consumer started: queue %s, stops on Ctrl+C or after -idle=%s without a message\n", targetTopic, *idle)
 
	var idleTimer *time.Timer
	var idleC <-chan time.Time
	if *idle > 0 {
		idleTimer = time.NewTimer(*idle)
		idleC = idleTimer.C
	}
 
	delivered := map[string]int{} // submission_id -> number of deliveries
	deliveries := 0
 
loop:
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				break loop
			}
			id := msg.Metadata.Get("submission_id")
			deliveries++
			delivered[id]++
			fmt.Fprintf(out, "received: submission_id=%s message_uuid=%s delivery=%d\n", id, msg.UUID, delivered[id])
			msg.Ack()
			if idleTimer != nil {
				idleTimer.Reset(*idle)
			}
		case <-idleC:
			fmt.Fprintf(out, "no message for %s, stopping\n", *idle)
			break loop
		case <-ctx.Done():
			break loop
		}
	}
 
	// What is really stored.
	inDB := map[string]bool{}
	rows, err := db.Query(`SELECT id FROM spk03_submissions`)
	must(err)
	for rows.Next() {
		var id string
		must(rows.Scan(&id))
		inDB[id] = true
	}
	must(rows.Err())
	rows.Close()
 
	dataWithoutMessage, messageWithoutData := 0, 0
	for id := range inDB {
		if delivered[id] == 0 {
			dataWithoutMessage++
			fmt.Fprintf(out, "DATA WITHOUT MESSAGE: submission_id=%s\n", id)
		}
	}
	for id := range delivered {
		if !inDB[id] {
			messageWithoutData++
			fmt.Fprintf(out, "MESSAGE WITHOUT DATA: submission_id=%s\n", id)
		}
	}
 
	fmt.Fprintf(out, "SUMMARY deliveries=%d unique=%d duplicates=%d db_rows=%d data_without_message=%d message_without_data=%d\n",
		deliveries, len(delivered), deliveries-len(delivered), len(inDB), dataWithoutMessage, messageWithoutData)
}
 
func must(err error) {
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
		panic(err)
	}
}
 