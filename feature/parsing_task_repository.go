package feature

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ParsingTaskRepository struct {
	session *mgo.Session
}

// Create new repository
func NewParsingTaskRepository(session *mgo.Session) ParsingTaskRepository {
	return ParsingTaskRepository{
		session: session,
	}
}

// Get collection of repository
func (p *ParsingTaskRepository) getCollection() *mgo.Collection {
	return p.session.DB("kadastr").C("feature_parsing_task")
}

// Find all tasks
func (p *ParsingTaskRepository) FindAll() []ParsingTask {
	var tasks []ParsingTask

	p.getCollection().Find(bson.M{}).All(&tasks)

	return tasks
}

// Insert task
func (p *ParsingTaskRepository) Insert(task ParsingTask) error {
	return p.getCollection().Insert(task)
}

// Update task
func (p *ParsingTaskRepository) Update(task ParsingTask) error {
	return p.getCollection().UpdateId(task.ID, task)
}

// Find by id
func (p *ParsingTaskRepository) FindById(id bson.ObjectId) (ParsingTask, error) {
	var task ParsingTask
	err := p.getCollection().Find(bson.M{"_id": id}).One(&task)

	return task, err
}

// Find one task by quarter
func (p *ParsingTaskRepository) FindByQuarter(quarter string) (ParsingTask, error) {
	var task ParsingTask
	err := p.getCollection().Find(bson.M{"quarter": quarter}).One(&task)

	return task, err
}
