package feature

import "gopkg.in/mgo.v2/bson"

type Entity struct {
	ID          bson.ObjectId `json:"oid" bson:"_id"`
	Address     string        `json:"address,omitempty" bson:"address"`
	AreaType    string        `json:"area_type,omitempty" bson:"area_type"`
	AreaUnit    string        `json:"area_unit,omitempty" bson:"area_unit"`
	AreaValue   float64       `json:"area_value,omitempty" bson:"area_value"`
	CadNumber   string        `json:"cn,omitempty" bson:"cad_number"`
	CadQuarter  string        `json:"kvartal_cn,omitempty" bson:"cad_quarter"`
	CadDistrict string        `json:"rayon_cn,omitempty" bson:"cad_district"`
	CadRegion   string        `json:"okrug_cn,omitempty" bson:"cad_region"`
}
