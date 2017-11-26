package api_server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/iamsalnikov/kadastr/parser"
)

type Server struct {
	Mongo *mgo.Session
	AMQP  *amqp.Connection
	Addr  string
}

func (s *Server) Run() {
	r := mux.NewRouter()

	featureTaskSender, err := parser.NewFeatureTaskSender(s.AMQP)
	if err != nil {
		log.Fatalf("Не удалось создать отправщика в очередь парсера: %s\n", err.Error())
	}

	quarterCheckSender, err := parser.NewQuarterCheckSender(s.AMQP)
	if err != nil {
		log.Fatalf("Не удалось создать отправщика в очередь квартального проверяльщика: %s\n", err.Error())
	}

	featureController := NewFeatureController(s.Mongo, featureTaskSender, quarterCheckSender)

	r.HandleFunc("/list-parsing", featureController.GetListParsing).Methods(http.MethodGet)
	r.HandleFunc("/add-parsing", featureController.AddParsingTask).Methods(http.MethodPost)
	r.HandleFunc("/search", featureController.FindFeature).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         s.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
