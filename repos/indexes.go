package repos

import (
	"gopkg.in/mgo.v2"
	"log"
)

// Create database indexes
func CreateIndexes(session *mgo.Session) {
	createFeatureIndex(session)
	createParsingTaskIndexes(session)
}

func createFeatureIndex(session *mgo.Session) {
	c := session.DB("kadastr").C("feature")

	if !isIndexExists(c, "quarter_area") {
		log.Println("Create index")
		c.EnsureIndex(mgo.Index{
			Key:  []string{"cad_quarter", "area_value"},
			Name: "quarter_area",
		})
	}

	if !isIndexExists(c, "cad_number") {
		log.Println("Create index cad_number")
		c.EnsureIndex(mgo.Index{
			Key:  []string{"cad_number"},
			Name: "cad_number",
		})
	}
}

func createParsingTaskIndexes(session *mgo.Session) {
	c := session.DB("kadastr").C("feature_parsing_task")

	if !isIndexExists(c, "quarter") {
		log.Println("Create index \"quarter\" for collection \"feature_parsing_task\"")
		c.EnsureIndex(mgo.Index{
			Key:  []string{"quarter"},
			Name: "quarter",
		})
	}
}

func isIndexExists(c *mgo.Collection, name string) bool {
	indexes, _ := c.Indexes()

	for _, index := range indexes {
		if index.Name == name {
			return true
		}
	}

	return false
}
