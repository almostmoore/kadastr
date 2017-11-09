package parser

import (
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
)

type FeatureTaskSender struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewFeatureTaskSender(conn *amqp.Connection) (*FeatureTaskSender, error) {
	fts := &FeatureTaskSender{
		conn: conn,
	}

	var err error
	fts.ch, err = fts.conn.Channel()
	if err != nil {
		return fts, err
	}

	fts.queue, err = fts.ch.QueueDeclare(
		FeatureParseQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	return fts, err
}

func (f *FeatureTaskSender) Send(taskId bson.ObjectId) error {
	return f.ch.Publish(
		"",
		FeatureParseQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(taskId.Hex()),
		},
	)
}
