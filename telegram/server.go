package telegram

import (
	"github.com/iamsalnikov/kadastr/parser"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type Server struct {
	APIToken      string
	Mongo         *mgo.Session
	AMQP          *amqp.Connection
	commandRouter *CommandRouter
}

func (s *Server) Run() {
	bot, err := tgbotapi.NewBotAPI(s.APIToken)
	if err != nil {
		log.Fatal("Не удалось запустить бота", err)
	}

	log.Printf("Авторизован под аккаунтом %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal("Не удалось подписаться на обновления", err)
	}

	s.BindCommands()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		s.ProcessMessage(bot, &update)
	}
}

func (s *Server) BindCommands() {
	s.commandRouter = &CommandRouter{}

	featureTaskSender, err := parser.NewFeatureTaskSender(s.AMQP)
	if err != nil {
		log.Fatalf("Не удалось создать отправщика в очередь парсера: %s\n", err.Error())
	}

	quarterCheckSender, err := parser.NewQuarterCheckSender(s.AMQP)
	if err != nil {
		log.Fatalf("Не удалось создать отправщика в очередь квартального проверяльщика: %s\n", err.Error())
	}

	featureProcessor := NewFeatureProcessor(s.Mongo, quarterCheckSender)
	s.commandRouter.AddProcessor("search", featureProcessor)
	s.commandRouter.SetDefaultCommand("search")

	addParsingTaskProcessor := NewAddParsingTaskProcessor(s.Mongo, featureTaskSender)
	s.commandRouter.AddProcessor("doparsing", addParsingTaskProcessor)

	listParsingTaskProcessor := NewListParsingTaskProcessor(s.Mongo)
	s.commandRouter.AddProcessor("listparsing", listParsingTaskProcessor)

	helpProcessor := HelpProcessor{}
	s.commandRouter.AddProcessor("start", helpProcessor)
	s.commandRouter.AddProcessor("help", helpProcessor)
}

func (s *Server) ProcessMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	answer, err := s.commandRouter.Run(update)

	if err != nil {
		answer = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка. "+err.Error())
	}

	answer.ParseMode = tgbotapi.ModeMarkdown

	bot.Send(answer)
}
