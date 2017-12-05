package models

import (
	"gopkg.in/mgo.v2/bson"
)

//GameRoom represents the room of players
type GameRoom struct {
	Players   []*TypieBird
	Available bool
}

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

//Updates represents updates that can be made to a typie bird
type Updates struct {
	Record float32 `json:"record"`
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
	bird.Record = updates.Record
	return nil
}
