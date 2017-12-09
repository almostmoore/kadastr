package parser

import (
	"github.com/almostmoore/kadastr/feature"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"log"
)

const QuarterCheckerQueue = "quarter_checker"

type QuarterCheckerCmd struct {
	Conn     *amqp.Connection
	Session  *mgo.Session
	ch       *amqp.Channel
	sender   *FeatureTaskSender
	taskRepo feature.ParsingTaskRepository
}

func (f *QuarterCheckerCmd) Run() {
	var err error
	f.ch, err = f.Conn.Channel()
	if err != nil {
		log.Fatal("Не удалось открыть канал", err)
	}

	defer f.ch.Close()

	q, err := f.ch.QueueDeclare(
		QuarterCheckerQueue,
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

	f.sender, err = NewFeatureTaskSender(f.Conn)
	if err != nil {
		log.Fatalf("Не удалось создать отправителя: %s\n", err.Error())
	}

	f.taskRepo = feature.NewParsingTaskRepository(f.Session)

	for {
		for m := range messages {
			f.processMessage(m)
		}
	}
}

func (f *QuarterCheckerCmd) processMessage(message amqp.Delivery) {
	cadQuarter := feature.ClearLeadZero(string(message.Body))
	log.Println("Получил задание на проверку парсинга квартала", cadQuarter)

	task, err := EnsureNewParsingTask(&f.taskRepo, cadQuarter)
	if err != nil {
		log.Printf("Ошибка: %s\n", err.Error())
	} else {
		err = f.sender.Send(task.ID)
		if err != nil {
			log.Printf("Произошла ошибка при отправлении квартала %s на парсинг: %s\n", cadQuarter, err.Error())
		} else {
			log.Printf("Квартал %s отправлен на парсинг\n", cadQuarter)
		}
	}

	f.ch.Ack(message.DeliveryTag, false)
	log.Printf("Проверка квартала %s завершена\n", cadQuarter)
}
