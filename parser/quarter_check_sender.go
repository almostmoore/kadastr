package parser

import (
	"github.com/streadway/amqp"
)

type QuarterCheckSender struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewQuarterCheckSender(conn *amqp.Connection) (*QuarterCheckSender, error) {
	fts := &QuarterCheckSender{
		conn: conn,
	}

	var err error
	fts.ch, err = fts.conn.Channel()
	if err != nil {
		return fts, err
	}

	fts.queue, err = fts.ch.QueueDeclare(
		QuarterCheckerQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	return fts, err
}

func (f *QuarterCheckSender) Send(quarter string) error {
	return f.ch.Publish(
		"",
		QuarterCheckerQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(quarter),
		},
	)
}
