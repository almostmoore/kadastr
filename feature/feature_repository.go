package feature

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FeatureRepository struct {
	session *mgo.Session
}

// NewFeature is constructor for FeatureRepository
func NewFeatureRepository(session *mgo.Session) FeatureRepository {
	return FeatureRepository{
		session: session,
	}
}

func (f *FeatureRepository) getCollection() *mgo.Collection {
	return f.session.DB("kadastr").C("feature")
}

func (f *FeatureRepository) FindAllByQuarterAndArea(quarter string, area float64) []Entity {
	var features []Entity

	f.getCollection().Find(bson.M{
		"cad_quarter": quarter,
		"area_value":  area,
	}).All(&features)

	return features
}

func (f *FeatureRepository) FindByCadNumber(number string) (Entity, error) {
	var entity Entity

	err := f.getCollection().Find(bson.M{
		"cad_number": number,
	}).One(&entity)

	return entity, err
}

func (f *FeatureRepository) Insert(entity Entity) error {
	return f.getCollection().Insert(entity)
}
