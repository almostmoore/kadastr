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
)

type Server struct {
	Mongo *mgo.Session
	AMQP  *amqp.Connection
	Addr  string
}

func (s *Server) Run() {
	r := mux.NewRouter()

	featureController := NewFeatureController(s.Mongo)

	r.HandleFunc("/listparsing", featureController.GetListParsing)

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         s.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
