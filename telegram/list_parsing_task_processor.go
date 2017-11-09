package telegram

import (
	"bytes"
	"fmt"
	"github.com/iamsalnikov/kadastr/feature"
	"gopkg.in/mgo.v2"
	"gopkg.in/telegram-bot-api.v4"
)

type ListParsingTaskProcessor struct {
	session *mgo.Session
}

func NewListParsingTaskProcessor(session *mgo.Session) ListParsingTaskProcessor {
	return ListParsingTaskProcessor{
		session: session,
	}
}

func (lptp ListParsingTaskProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	repo := feature.NewParsingTaskRepository(lptp.session)
	parsingTasks := repo.FindAll()

	answer := bytes.NewBufferString("Кварталы на парсинге:\n")

	for _, task := range parsingTasks {
		answer.WriteString(fmt.Sprintf("\n*%s* - %s", task.Quarter, task.TextStatus))
	}

	return tgbotapi.NewMessage(upd.Message.Chat.ID, answer.String()), nil
}
