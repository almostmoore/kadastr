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

func main() {
	flag.StringVar(&mongoConnectionString, "mongo", os.Getenv("MONGO"), "")
	flag.StringVar(&rabbitConnectionString, "rabbit", os.Getenv("RABBIT"), "")
	flag.StringVar(&tgToken, "tgtoken", os.Getenv("TG_TOKEN"), "")
	flag.StringVar(&mode, "mode", "tg", "")
	flag.StringVar(&apiAddr, "addr", os.Getenv("ADDR"), "0.0.0.0:8080")
	flag.Parse()

	session, err := mgo.Dial("mongodb://" + mongoConnectionString)
	if err != nil {
		log.Fatal("Не удалось получить доступ в БД", err)
	}
	err = session.Ping()
	if err != nil {
		log.Fatal("Не удалось соединиться с БД", err)
	}

	repos.CreateIndexes(session)

	// Rabbit connection
	rConn, err := amqp.Dial("amqp://" + rabbitConnectionString)
	if err != nil {
		log.Fatal("Не удалось соединиться с раббитом", err)
	}

	go func() {
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
	}()

	// Feature parser worker
	if mode == "rb" {
		worker := &parser.FeatureParserCmd{
			Conn:    rConn,
			Session: session,
		}

		worker.Run()
	}

	// Quarter checker worker
	if mode == "quarter" {
		worker := &parser.QuarterCheckerCmd{
			Conn:    rConn,
			Session: session,
		}

		worker.Run()
	}

	// Telegram client
	if mode == "tg" {
		tg := &telegram.Server{
			APIToken: tgToken,
			Mongo:    session,
			AMQP:     rConn,
		}

		tg.Run()
	}

	// Api server
	if mode == "api" {
		api := &api_server.Server{
			Mongo: session,
			AMQP:  rConn,
			Addr:  apiAddr,
		}

		api.Run()
	}
}
