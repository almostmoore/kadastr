package feature

import "gopkg.in/mgo.v2/bson"

const (
	ParsingStatusPrepare  = 1
	ParsingStatusProgress = 2
	ParsingStatusReady    = 3
)

type ParsingTask struct {
	ID         bson.ObjectId `json:"oid" bson:"_id"`
	Quarter    string        `json:"quarter" bson:"quarter"`
	TextStatus string        `json:"text_status" bson:"text_status"`
	Status     int64         `json:"status" bson:"status"`
}
