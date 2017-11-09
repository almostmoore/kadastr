package parser

import (
	"github.com/iamsalnikov/kadastr/feature"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const FeatureParseQueue = "feature_parsing_task"

type FeatureParserCmd struct {
	Conn     *amqp.Connection
	Session  *mgo.Session
	ch       *amqp.Channel
	taskRepo feature.ParsingTaskRepository
}

func (f *FeatureParserCmd) Run() {
	var err error
	f.ch, err = f.Conn.Channel()
	if err != nil {
		log.Fatal("Не удалось открыть канал", err)
	}

	defer f.ch.Close()

	q, err := f.ch.QueueDeclare(
		FeatureParseQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Не удалось открыть очередь", err)
	}

	messages, err := f.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Не удалось подписаться на очередь")
	}

	f.taskRepo = feature.NewParsingTaskRepository(f.Session)

	for {
		for m := range messages {
			f.processMessage(m)
		}
	}
}

func (f *FeatureParserCmd) processMessage(message amqp.Delivery) {
	taskId := string(message.Body)
	log.Println("Получил задание на парсинг", taskId)

	task, err := f.taskRepo.FindById(bson.ObjectIdHex(taskId))
	if err != nil {
		log.Println("Задача не найдена", err)
		f.ch.Ack(message.DeliveryTag, false)
		return
	}

	task.Status = feature.ParsingStatusProgress
	task.TextStatus = "В обработке"
	f.taskRepo.Update(task)

	log.Println("Парсинг квартала ", task.Quarter)

	parser := NewFeatureParser(f.Session)
	parser.Run(task.Quarter, 3)

	log.Println("Парсинг квартала окончен ", task.Quarter)
	f.ch.Ack(message.DeliveryTag, false)

	task.Status = feature.ParsingStatusReady
	task.TextStatus = "Обработано"
	f.taskRepo.Update(task)
}
