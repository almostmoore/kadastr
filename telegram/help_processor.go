package telegram

import "gopkg.in/telegram-bot-api.v4"

type HelpProcessor struct{}

func (h HelpProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	return tgbotapi.NewMessage(upd.Message.Chat.ID, `
*Кадастробот*

Бот выполняет поиск кадастровых номеров по кадастровому кварталу и площади.

*Примеры использования*

1. Поиск участка с площадью 1500 кв. м. в кадастровом квартале 29:08:104222

Отправить сообщение: 29:08:104222 1500

2. Добавление кадастрового квартала на парсинг:

/doparsing 29:08:104222

3. Добавление нескольких кадастровых кварталов на парсинг:

/doparsing 29:08:104222 29:08:104223 29:08:104224

4. Статус парсинга кадастровых кварталов:

/listparsing
	`), nil
}
