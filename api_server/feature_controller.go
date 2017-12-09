package api_server

import (
	"encoding/json"
	"github.com/almostmoore/kadastr/feature"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/almostmoore/kadastr/parser"
	"github.com/almostmoore/kadastr/api_server/messages"
	"log"
	"strconv"
)

type FeatureController struct {
	session           *mgo.Session
	featureTaskSender *parser.FeatureTaskSender
	quarterCheckSender *parser.QuarterCheckSender
}

func NewFeatureController(session *mgo.Session, fts *parser.FeatureTaskSender, qcs *parser.QuarterCheckSender) FeatureController {
	return FeatureController{
		session:            session,
		featureTaskSender:  fts,
		quarterCheckSender: qcs,
	}
}

func (f *FeatureController) GetListParsing(resp http.ResponseWriter, req *http.Request) {
	repo := feature.NewParsingTaskRepository(f.session)
	tasks := repo.FindAll()

	responseTasks := make([]messages.ParsingTask, 0, len(tasks))
	for _, task := range tasks {
		responseTasks = append(responseTasks, messages.ParsingTask{
			Quarter: task.Quarter,
			TextStatus: task.TextStatus,
			Status: task.Status,
		})
	}

	resp.WriteHeader(http.StatusOK)
	resp.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(resp)
	encoder.Encode(responseTasks)
}

func (f *FeatureController) AddParsingTask(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	quarters := make([]string, 0)

	resp.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&quarters)
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
			err = f.featureTaskSender.Send(task.ID)
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

	encoder := json.NewEncoder(resp)
	encoder.Encode(answer)
}

func (f *FeatureController) FindFeature(resp http.ResponseWriter, req *http.Request) {
	quarter := req.URL.Query().Get("quarter")
	square := req.URL.Query().Get("square")

	repo := feature.NewFeatureRepository(f.session)
	s, _ := strconv.ParseFloat(square, 64)

	encoder := json.NewEncoder(resp)
	resp.Header().Add("Content-Type", "application/json")

	if quarter == "" || square == "" {
		err := messages.Error{
			Message: "Нужно указать квартал и площадь. Например - 29:08:103701 2578",
		}

		resp.WriteHeader(http.StatusBadRequest)
		encoder.Encode(err)
		return
	}

	f.quarterCheckSender.Send(quarter)

	features := repo.FindAllByQuarterAndArea(quarter, s)
	foundedFeatures := make([]messages.FindFeature, 0, len(features))

	for _, ft := range features {
		foundedFeatures = append(foundedFeatures, messages.FindFeature{
			CadNumber: ft.CadNumber,
			Address: ft.Address,
		})
	}

	encoder.Encode(foundedFeatures)
}
