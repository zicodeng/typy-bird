package models

import (
	"gopkg.in/mgo.v2/bson"
)

//TypieBird represents a player
type TypieBird struct {
	UserName string        `json:"userName"`
	Record   float32       `json:"record"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	Record float32 `json:"record"`
}

//UpdateRecord updates a typie bird's score to the highest score
func (bird *TypieBird) UpdateRecord(updates *Updates) error {
	if updates.Record < bird.Record {
		bird.Record = updates.Record
	}
	return nil
}