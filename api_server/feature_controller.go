package api_server

import (
	"encoding/json"
	"fmt"
	"github.com/iamsalnikov/kadastr/feature"
	"gopkg.in/mgo.v2"
	"net/http"
)

type FeatureController struct {
	session *mgo.Session
}

func NewFeatureController(session *mgo.Session) FeatureController {
	return FeatureController{
		session: session,
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
