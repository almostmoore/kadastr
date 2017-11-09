package telegram

import "gopkg.in/telegram-bot-api.v4"

type CommandRouter struct {
	processors     map[string]CommandProcessor
	defaultCommand string
}

func (cr *CommandRouter) AddProcessor(command string, processor CommandProcessor) {
	if cr.processors == nil {
		cr.processors = make(map[string]CommandProcessor)
	}

	cr.processors[command] = processor
}

func (cr *CommandRouter) SetDefaultCommand(command string) {
	cr.defaultCommand = command
}

func (cr *CommandRouter) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	processAsCommand := upd.Message.IsCommand() || cr.defaultCommand != ""
	if !processAsCommand {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, ""), UndefinedCommand
	}

	processor, ok := cr.processors[upd.Message.Command()]
	if ok {
		return processor.Run(upd)
	}

	processor, ok = cr.processors[cr.defaultCommand]
	if ok {
		return processor.Run(upd)
	}

	return tgbotapi.NewMessage(upd.Message.Chat.ID, "Неизвестная команда"), UndefinedCommand
}
