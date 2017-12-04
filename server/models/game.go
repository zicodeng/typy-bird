package models 

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


type TypieBird struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	UserName  string        `json:"userName"`
	Record time.Time
}

type Updates struct {
	Record	time.Time
}
