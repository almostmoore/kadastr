package telegram

import "gopkg.in/telegram-bot-api.v4"

type CommandProcessor interface {
	Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error)
}
