package game 

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type TypieBird struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	UserName  string        `json:"userName"`
	Record		time.Time			`json:"record"`
}

func (bird *TypieBird) UpdateRecord(record time.Time) error {
	if(record < bird.Record) {
		bird.Record = record;
	}
	return nil
}
