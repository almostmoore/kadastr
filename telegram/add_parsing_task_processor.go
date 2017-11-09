package telegram

import (
	"bytes"
	"fmt"
	"github.com/iamsalnikov/kadastr/feature"
	"github.com/iamsalnikov/kadastr/parser"
	"gopkg.in/mgo.v2"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
)

type AddParsingTaskProcessor struct {
	session           *mgo.Session
	quarterRegexp     *regexp.Regexp
	leadingZeroRegexp *regexp.Regexp
	sender            *parser.FeatureTaskSender
}

func NewAddParsingTaskProcessor(session *mgo.Session, sender *parser.FeatureTaskSender) AddParsingTaskProcessor {
	return AddParsingTaskProcessor{
		session:           session,
		quarterRegexp:     regexp.MustCompile("(\\d+?:\\d+?:\\d+)\\s*"),
		leadingZeroRegexp: regexp.MustCompile("(:)(0+)(\\d+)"),
		sender:            sender,
	}
}

func (aptp AddParsingTaskProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	quarters := aptp.extractQuarters(upd.Message.CommandArguments())
	repo := feature.NewParsingTaskRepository(aptp.session)

	answer := bytes.NewBufferString("Реузльтат добавления кварталов на парсинг:\n")

	for _, quarter := range quarters {
		task, sendToParsing, err := parser.EnsureParsingTask(&repo, quarter)

		if err != nil {
			answer.WriteString(fmt.Sprintf("\n*%s* - не добавлен (%s)", quarter, err.Error()))
		} else if sendToParsing {
			err = aptp.sender.Send(task.ID)
			if err != nil {
				answer.WriteString(fmt.Sprintf("\n*%s* - не отправлен на парсинг (%s)", quarter, err.Error()))
			} else {
				answer.WriteString(fmt.Sprintf("\n*%s* - добавлен", quarter))
			}
		} else {
			answer.WriteString(fmt.Sprintf("\n*%s* - добавлен", quarter))
		}
	}

	return tgbotapi.NewMessage(upd.Message.Chat.ID, answer.String()), nil
}

func (aptp AddParsingTaskProcessor) extractQuarters(str string) []string {
	sub := aptp.quarterRegexp.FindAllStringSubmatch(str, -1)

	data := make([]string, 0, len(sub))
	for _, s := range sub {
		data = append(data, feature.ClearLeadZero(s[1]))
	}

	return data
}
