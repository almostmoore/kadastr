package quarter

import "gopkg.in/mgo.v2/bson"

type Entity struct {
	ID                bson.ObjectId `bson:"_id" json:"id"`
	Quarter           string        `bson:"quarter" json:"quarter"`
	Status            string        `bson:"status" json:"status"`
	ProcessedElements uint64        `bson:"processed_elements" json:"processed_elements"`
}
