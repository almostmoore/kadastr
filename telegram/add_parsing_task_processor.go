package telegram

import (
	"bytes"
	"fmt"
	"github.com/almostmoore/kadastr/feature"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"github.com/almostmoore/kadastr/api_server"
	"strings"
)

type AddParsingTaskProcessor struct {
	ApiClient         *api_server.Client
	quarterRegexp     *regexp.Regexp
	leadingZeroRegexp *regexp.Regexp
}

func NewAddParsingTaskProcessor(apiClient *api_server.Client) AddParsingTaskProcessor {
	return AddParsingTaskProcessor{
		ApiClient:         apiClient,
		quarterRegexp:     regexp.MustCompile("(\\d+?:\\d+?:\\d+)\\s*"),
		leadingZeroRegexp: regexp.MustCompile("(:)(0+)(\\d+)"),
	}
}

func (aptp AddParsingTaskProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	quarters := aptp.extractQuarters(upd.Message.CommandArguments())
	addResult, err := aptp.ApiClient.AddParsingTask(quarters)
	if err != nil {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, "При обращении к серверу произошла ошибка\n"), nil
	}

	answer := bytes.NewBufferString("Реузльтат добавления кварталов на парсинг:\n")

	answer.WriteString(fmt.Sprintf("*Добавлено на парсинг - %d*", len(addResult.Added)))
	if len(addResult.Added) > 0 {
		answer.WriteString(": " + strings.Join(addResult.Added, ", "))
	}

	answer.WriteString(fmt.Sprintf("\n*Не добавлено на парсинг - %d*", len(addResult.NotAdded)))
	if len(addResult.NotAdded) > 0 {
		answer.WriteString(": " + strings.Join(addResult.Added, ", "))
	}

	answer.WriteString(fmt.Sprintf("\n*Не отправлено на парсинг - %d*", len(addResult.NotSent)))
	if len(addResult.NotSent) > 0 {
		answer.WriteString(": " + strings.Join(addResult.NotSent, ", "))
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
