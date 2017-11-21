package api_server

import (
	"encoding/json"
	"fmt"
	"github.com/iamsalnikov/kadastr/feature"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/iamsalnikov/kadastr/parser"
	"github.com/iamsalnikov/kadastr/api_server/messages"
	"log"
	"bytes"
)

type FeatureController struct {
	session *mgo.Session
	sender *parser.FeatureTaskSender
}

func NewFeatureController(session *mgo.Session, sender *parser.FeatureTaskSender) FeatureController {
	return FeatureController{
		session: session,
		sender: sender,
	}
}

func (f *FeatureController) GetListParsing(resp http.ResponseWriter, req *http.Request) {
	repo := feature.NewParsingTaskRepository(f.session)
	parsingTasks := repo.FindAll()

	resp.WriteHeader(http.StatusOK)
	resp.Header().Add("Content-Type", "application/json")
	answer, _ := json.Marshal(parsingTasks)

	fmt.Fprintln(resp, string(answer))
}

func (f *FeatureController) AddParsingTask(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	quarters := make([]string, 0)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&quarters)

	resp.Header().Add("Content-Type", "application/json")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}

	repo := feature.NewParsingTaskRepository(f.session)

	answer := &messages.AddParsingTaskAnswer{}
	for _, quarter := range quarters {
		quarter = feature.ClearLeadZero(quarter)
		task, sendToParsing, err := parser.EnsureParsingTask(&repo, quarter)

		if err != nil {
			answer.NotAdded = append(answer.NotAdded, quarter)
			log.Printf("Квартал %s - не добавлен (%s)\n", quarter, err.Error())
		} else if sendToParsing {
			err = f.sender.Send(task.ID)
			if err != nil {
				answer.NotSent = append(answer.NotSent, quarter)
				log.Printf("Квартал %s - не отправлен на парсинг (%s)\n", quarter, err.Error())
			} else {
				answer.Added = append(answer.Added, quarter)
				log.Printf("Квартал %s - добавлен\n", quarter)
			}
		} else {
			answer.Added = append(answer.Added, quarter)
			log.Printf("Квартал %s - добавлен\n", quarter)
		}
	}

	jsonAnswer := bytes.NewBufferString("")
	encoder := json.NewEncoder(jsonAnswer)
	encoder.Encode(answer)

	fmt.Fprintln(resp, jsonAnswer.String())
}
