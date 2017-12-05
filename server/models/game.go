package models

import (
	"gopkg.in/mgo.v2/bson"
)

//NewTypieBird represents a creating a new player
type NewTypieBird struct {
	UserName string `json:"userName"`
}

//TypieBird represents a player
type TypieBird struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	UserName string        `json:"userName"`
	Record   float32       `json:"record"`
	Position int           `json:"position"`
}

//Credentials represents typie bird credentials required to join a game
type Credentials struct {
	UserName string `json:"userName"`
}

//Updates represents updates that can be made to a typie bird
type Updates struct {
	Record   float32 `json:"record"`
	Position int     `json:"position"`
}

//ToTypie takes a NewTypie and turns it into a Typie
func (nt *NewTypieBird) ToTypie() *TypieBird {
	return &TypieBird{
		ID:       bson.NewObjectId(),
		UserName: nt.UserName,
		Record:   0,
		Position: 0,
	}
}

//Update updates a typie bird's position and score
func (bird *TypieBird) Update(updates *Updates) error {
	if updates.Record < bird.Record {
		bird.Record = updates.Record
	}
	bird.Position = updates.Position
	return nil
}
