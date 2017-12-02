package game 

import (
	"gopkg.in/mgo.v2/bson"

)


type TypieBird struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	UserName  string        `json:"userName"`
}