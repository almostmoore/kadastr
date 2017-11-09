package quarter

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	session *mgo.Session
}

func NewQuarterRepository(session *mgo.Session) Repository {
	return Repository{
		session: session,
	}
}

func (r *Repository) getCollection() *mgo.Collection {
	return r.session.DB("kadastr").C("quarter")
}

func (r *Repository) Insert(entity Entity) error {
	return r.getCollection().Insert(entity)
}

func (r *Repository) Update(id bson.ObjectId, entity Entity) error {
	return r.getCollection().Update(bson.M{"_id": id}, entity)
}

func (r *Repository) Delete(id bson.ObjectId) error {
	return r.getCollection().RemoveId(id)
}

func (r *Repository) FindAll() []Entity {
	var list []Entity

	r.getCollection().Find(bson.M{}).All(&list)

	return list
}
