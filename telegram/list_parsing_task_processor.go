package telegram

import (
	"bytes"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/iamsalnikov/kadastr/api_server"
)

type ListParsingTaskProcessor struct {
	ApiClient *api_server.Client
}

func NewListParsingTaskProcessor(apiClient *api_server.Client) ListParsingTaskProcessor {
	return ListParsingTaskProcessor{
		ApiClient: apiClient,
	}
}

func (lptp ListParsingTaskProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	parsingTasks, err := lptp.ApiClient.GetParsingTasksList()
	if err != nil {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, "При обращении к серверу произошла ошибка\n"), nil
	}

	answer := bytes.NewBufferString("Кварталы на парсинге:\n")

	for _, task := range parsingTasks {
		answer.WriteString(fmt.Sprintf("\n*%s* - %s", task.Quarter, task.TextStatus))
	}

	return tgbotapi.NewMessage(upd.Message.Chat.ID, answer.String()), nil
}
