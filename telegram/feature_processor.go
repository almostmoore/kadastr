package telegram

import (
	"bytes"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strings"
	"github.com/almostmoore/kadastr/api_server"
)

type FeatureProcessor struct {
	ApiClient *api_server.Client
	qsRegExp  *regexp.Regexp
}

func NewFeatureProcessor(apiClient *api_server.Client) FeatureProcessor {
	return FeatureProcessor{
		ApiClient: apiClient,
		qsRegExp:  regexp.MustCompile("\\s*(\\d+?:\\d+?:\\d+)\\s+(\\d+[.,]?\\d*)"),
	}
}

func (f FeatureProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	quarter, square := f.extractQuarterAndSquare(upd.Message.Text)

	features, searchError, err := f.ApiClient.FindFeature(quarter, square)
	if err != nil {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, "Произошла ошибка при обращении к серверу\n"), nil
	}

	if searchError.Message != "" {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, searchError.Message), nil
	}

	answer := bytes.NewBufferString(fmt.Sprintf("Поиск для кадастрового квартала *%s* и площади *%s*\n", quarter, square))
	for _, f := range features {
		answer.WriteString(fmt.Sprintf("\n*%s* - %s", f.CadNumber, f.Address))
	}

	if len(features) == 0 {
		answer.WriteString("\nНичего не найдено")
	}

	return tgbotapi.NewMessage(upd.Message.Chat.ID, answer.String()), nil
}

func (f FeatureProcessor) extractQuarterAndSquare(str string) (string, string) {
	found := f.qsRegExp.FindStringSubmatch(str)

	if len(found) != 3 {
		return "", ""
	}

	return found[1], strings.Replace(found[2], ",", ".", -1)
}
