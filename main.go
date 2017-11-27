package main

import (
	"flag"
	"github.com/iamsalnikov/kadastr/api_server"
	"github.com/iamsalnikov/kadastr/parser"
	"github.com/iamsalnikov/kadastr/repos"
	"github.com/iamsalnikov/kadastr/telegram"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

var mongoConnectionString string
var rabbitConnectionString string
var mode string
var tgToken string
var apiAddr string
var session *mgo.Session
var rabbitConnection *amqp.Connection

func main() {
	flag.StringVar(&mongoConnectionString, "mongo", os.Getenv("MONGO"), "")
	flag.StringVar(&rabbitConnectionString, "rabbit", os.Getenv("RABBIT"), "")
	flag.StringVar(&tgToken, "tgtoken", os.Getenv("TG_TOKEN"), "")
	flag.StringVar(&mode, "mode", "tg", "")
	flag.StringVar(&apiAddr, "addr", os.Getenv("ADDR"), "0.0.0.0:8080")
	flag.Parse()

	initEnvironmentRoutine(mode)
	if mode != "tg" {
		go dbHealthCheck()
	}

	// Feature parser worker
	if mode == "rb" {
		worker := &parser.FeatureParserCmd{
			Conn:    rabbitConnection,
			Session: session,
		}

		worker.Run()
	}

	// Quarter checker worker
	if mode == "quarter" {
		worker := &parser.QuarterCheckerCmd{
			Conn:    rabbitConnection,
			Session: session,
		}

		worker.Run()
	}

	// Telegram client
	if mode == "tg" {
		tg := &telegram.Server{
			APIToken: tgToken,
			ApiClient: api_server.NewClient("http://" + apiAddr),
		}

		tg.Run()
	}

	// Api server
	if mode == "api" {
		api := &api_server.Server{
			Mongo: session,
			AMQP:  rabbitConnection,
			Addr:  apiAddr,
		}

		api.Run()
	}
}

func initEnvironmentRoutine(environment string) {
	if environment == "tg" {
		return
	}

	var err error
	session, err = mgo.Dial("mongodb://" + mongoConnectionString)
	if err != nil {
		log.Fatal("Не удалось получить доступ в БД", err)
	}
	err = session.Ping()
	if err != nil {
		log.Fatal("Не удалось соединиться с БД", err)
	}

	repos.CreateIndexes(session)

	// Rabbit connection
	rabbitConnection, err = amqp.Dial("amqp://" + rabbitConnectionString)
	if err != nil {
		log.Fatal("Не удалось соединиться с раббитом", err)
	}
}

func dbHealthCheck() {
	timer := time.Tick(10 * time.Second)
	for {
		select {
		case <-timer:
			err := session.Ping()
			if err != nil {
				log.Println("Can not ping DB", err.Error())
				session.Refresh()
			}
		}

	}
}