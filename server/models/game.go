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

//RecordUpdates represents updates to a bird's score
type RecordUpdates struct {
	Record float32 `json:"record"`
}

//PositionUpdates represents updates to a bird's position
type PositionUpdates struct {
	Position int `json:"position"`
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

//UpdateRecord updates a typie bird's position and score
func (bird *TypieBird) UpdateRecord(ru *RecordUpdates) error {
	if ru.Record < bird.Record {
		bird.Record = ru.Record
	}
	return nil
}

//UpdatePosition updates a typie bird's position and score
func (bird *TypieBird) UpdatePosition(pu *PositionUpdates) error {
	bird.Position = pu.Position
	return nil
}
