// SC-E publisher: the "application" side of the transactional outbox.
//
// For every submission it opens ONE database transaction that writes
//  1. the business row (table spk03_submissions), and
//  2. the event, as a row of the outbox table, through Watermill's forwarder decorator.
//
// The row is later moved to RabbitMQ by the forwarder program. The two writes are
// committed or rolled back together, which is the whole point of the pattern.
//
// Submission #3 (failAt) is the one that fails in the failure modes (-mode):
//
//	outbox                      nothing fails, 5 of 5 committed
//	outbox-rollback             #3: error after both writes, tx.Rollback()
//	outbox-crash-before-commit  #3: process killed (os.Exit) before Commit
//	outbox-crash-after-commit   #3: process killed right after Commit
//	direct-publish-first        BASELINE, no outbox: publish to RabbitMQ, then write the DB; #3's DB tx is rolled back
//	direct-commit-first         BASELINE, no outbox: Commit, then publish to RabbitMQ; #3 is killed between the two
//
// Run order: consumer first (it declares the queue), then forwarder, then this program.
// Every run starts with an empty spk03_submissions table. The outbox table is not touched.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"

	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/OUTBOX/outlog"
	"github.com/ITA-Internship/GoJudge/spike/spk-03/RabbitMQ/rmqenv"
)

const (
	forwarderTopic = "outbox_forwarder"   // outbox table: "watermill_outbox_forwarder"
	targetTopic    = "outbox_submissions" // RabbitMQ queue the events are meant for
	messages       = 5
	failAt         = 3

	createSubmissions = `CREATE TABLE IF NOT EXISTS spk03_submissions (
		id         UUID PRIMARY KEY,
		payload    TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`
)

var (
	out    io.Writer = os.Stdout
	logger watermill.LoggerAdapter
	schema = watermillSQL.DefaultPostgreSQLSchema{}
)

func main() {
	mode := flag.String("mode", "outbox", "outbox | outbox-rollback | outbox-crash-before-commit | outbox-crash-after-commit | direct-publish-first | direct-commit-first")
	logPath := flag.String("log", "", "also write the output to this log file")
	flag.Parse()

	var closeLog func()
	out, closeLog = outlog.Open(*logPath)
	defer closeLog()
	logger = watermill.NewStdLoggerWithOut(out, false, false)

	switch *mode {
	case "outbox", "outbox-rollback", "outbox-crash-before-commit", "outbox-crash-after-commit",
		"direct-publish-first", "direct-commit-first":
	default:
		panic("unknown -mode " + *mode)
	}
	ctx := context.Background()

	dsn, err := rmqenv.PostgresDSN() // credentials come from .env (POSTGRES_*)
	must(err)
	db, err := sql.Open("postgres", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	// The outbox table must exist before the first transaction uses it
	// (the publisher cannot create it inside a transaction, the forwarder creates it on start).
	for _, q := range schema.SchemaInitializingQueries(forwarderTopic) {
		_, err = db.Exec(q)
		must(err)
	}
	_, err = db.Exec(createSubmissions)
	must(err)
	_, err = db.Exec(`TRUNCATE spk03_submissions`)
	must(err)

	// The baseline modes publish straight to RabbitMQ, without the outbox.
	var direct *amqp.Publisher
	if strings.HasPrefix(*mode, "direct") {
		amqpURI, err := rmqenv.URI() // credentials come from .env (RABBITMQ_*)
		must(err)
		cfg := amqp.NewDurableQueueConfig(amqpURI)
		cfg.Publish.ConfirmDelivery = true // wait for the broker to confirm every message
		direct, err = amqp.NewPublisher(cfg, logger)
		must(err)
		defer direct.Close()
	}

	outboxTable := schema.MessagesTable(forwarderTopic)
	outboxBefore := count(db, outboxTable)
	committed, rolledBack := 0, 0

	for i := 1; i <= messages; i++ {
		id := watermill.NewUUID() // submission id, also the key of the message
		payload := fmt.Sprintf("submission #%d", i)
		msg := message.NewMessage(watermill.NewUUID(), []byte(payload))
		msg.Metadata.Set("submission_id", id)
		failing := i == failAt

		tx, err := db.BeginTx(ctx, nil)
		must(err)

		// 1. Business state.
		_, err = tx.ExecContext(ctx, `INSERT INTO spk03_submissions (id, payload) VALUES ($1, $2)`, id, payload)
		must(err)

		// 2. The event.
		if *mode == "direct-publish-first" {
			must(direct.Publish(targetTopic, msg)) // goes to the broker at once, outside of the DB transaction
			fmt.Fprintf(out, "#%d published directly to RabbitMQ: submission_id=%s\n", i, id)
			if failing {
				must(tx.Rollback())
				rolledBack++
				fmt.Fprintf(out, "#%d DB ROLLED BACK after the message was already published: submission_id=%s\n", i, id)
				continue
			}
		} else if strings.HasPrefix(*mode, "outbox") {
			must(publishToOutbox(tx, msg)) // INSERT into the outbox table, inside the same transaction
			if failing && *mode == "outbox-rollback" {
				must(tx.Rollback())
				rolledBack++
				fmt.Fprintf(out, "#%d ROLLED BACK (business row and outbox row): submission_id=%s\n", i, id)
				continue
			}
			if failing && *mode == "outbox-crash-before-commit" {
				crash(fmt.Sprintf("#%d killed before Commit: submission_id=%s", i, id))
			}
		}

		must(tx.Commit())
		committed++
		fmt.Fprintf(out, "#%d committed: submission_id=%s message_uuid=%s\n", i, id, msg.UUID)

		if failing && *mode == "outbox-crash-after-commit" {
			crash(fmt.Sprintf("#%d killed right after Commit: submission_id=%s", i, id))
		}
		if *mode == "direct-commit-first" {
			if failing {
				crash(fmt.Sprintf("#%d killed between Commit and Publish: submission_id=%s", i, id))
			}
			must(direct.Publish(targetTopic, msg))
			fmt.Fprintf(out, "#%d published directly to RabbitMQ: submission_id=%s\n", i, id)
		}
	}

	fmt.Fprintf(out, "SUMMARY mode=%s attempted=%d committed=%d rolled_back=%d submissions_in_db=%d outbox_rows_added=%d\n",
		*mode, messages, committed, rolledBack, count(db, "spk03_submissions"), count(db, outboxTable)-outboxBefore)
}

// publishToOutbox writes the event to the outbox table through the given transaction.
// forwarder.NewPublisher wraps the message in an envelope (the destination topic plus the
// original message) and publishes it to forwarderTopic. The forwarder program reads the
// envelope and publishes the original message to targetTopic in RabbitMQ.
func publishToOutbox(tx *sql.Tx, msg *message.Message) error {
	sqlPublisher, err := watermillSQL.NewPublisher(tx, watermillSQL.PublisherConfig{SchemaAdapter: schema}, logger)
	if err != nil {
		return err
	}
	outbox := forwarder.NewPublisher(sqlPublisher, forwarder.PublisherConfig{ForwarderTopic: forwarderTopic})
	defer outbox.Close()

	return outbox.Publish(targetTopic, msg)
}

func count(db *sql.DB, table string) int {
	var n int
	must(db.QueryRow("SELECT count(*) FROM " + table).Scan(&n))
	return n
}

// crash imitates a killed process: os.Exit skips every deferred call,
// and PostgreSQL rolls back the open transaction when the connection drops.
func crash(reason string) {
	fmt.Fprintln(out, "CRASH: "+reason)
	os.Exit(1)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(out, "ERROR:", err)
		panic(err)
	}
}
