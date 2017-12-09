package rapi

import "github.com/almostmoore/kadastr/feature"

type featureAnswer struct {
	Feature featureDataAnswer `json:"feature,omitempty"`
	Status  int32             `json:"status"`
}

type featureDataAnswer struct {
	Attributes feature.Entity `json:"attrs,omitempty"`
}
