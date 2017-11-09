package telegram

import (
	"bytes"
	"fmt"
	"github.com/iamsalnikov/kadastr/feature"
	"github.com/iamsalnikov/kadastr/parser"
	"gopkg.in/mgo.v2"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strconv"
	"strings"
)

type FeatureProcessor struct {
	session        *mgo.Session
	qsRegExp       *regexp.Regexp
	quarterChecker *parser.QuarterCheckSender
}

func NewFeatureProcessor(session *mgo.Session, quarterChecker *parser.QuarterCheckSender) FeatureProcessor {
	return FeatureProcessor{
		session:        session,
		qsRegExp:       regexp.MustCompile("\\s*(\\d+?:\\d+?:\\d+)\\s+(\\d+[.,]?\\d*)"),
		quarterChecker: quarterChecker,
	}
}

func (f FeatureProcessor) Run(upd *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	quarter, square := f.extractQuarterAndSquare(upd.Message.Text)

	repo := feature.NewFeatureRepository(f.session)
	s, _ := strconv.ParseFloat(square, 64)

	if quarter == "" || square == "" {
		return tgbotapi.NewMessage(upd.Message.Chat.ID, "Нужно указать квартал и площадь. Например - 29:08:103701 2578\n\n/help - справка"), nil
	}

	f.quarterChecker.Send(quarter)

	answer := bytes.NewBufferString(fmt.Sprintf("Поиск для кадастрового квартала *%s* и площади *%.2f*\n", quarter, s))
	features := repo.FindAllByQuarterAndArea(quarter, s)

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
