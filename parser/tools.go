package parser

import (
	"errors"
	"github.com/iamsalnikov/kadastr/feature"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ErrTaskAlreadyExists = errors.New("Task already exists")

func EnsureNewParsingTask(repo *feature.ParsingTaskRepository, quarter string) (feature.ParsingTask, error) {
	quarter = feature.ClearLeadZero(quarter)

	task, err := repo.FindByQuarter(quarter)
	if err == nil {
		return task, ErrTaskAlreadyExists
	} else if err != mgo.ErrNotFound {
		return task, err
	}

	task = feature.ParsingTask{
		ID:         bson.NewObjectId(),
		Quarter:    quarter,
		Status:     feature.ParsingStatusPrepare,
		TextStatus: "Приготовлен к парсингу",
	}

	return task, repo.Insert(task)
}

func EnsureParsingTask(repo *feature.ParsingTaskRepository, quarter string) (feature.ParsingTask, bool, error) {
	task, err := EnsureNewParsingTask(repo, quarter)
	if err == nil {
		return task, true, err
	} else if err == ErrTaskAlreadyExists && task.Status == feature.ParsingStatusReady {
		task.Status = feature.ParsingStatusPrepare
		task.TextStatus = "Приготовлен к парсингу"

		err = repo.Update(task)
		return task, true, err
	}

	return task, false, err
}
