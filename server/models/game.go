package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

//MAKE ERROR CONSTANT

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

//Update updates a typie bird's score
func (room *GameRoom) Update(typieBirdID bson.ObjectId, updates *Updates) (*TypieBird, error) {
	for _, player := range room.Players {
		if player.ID == typieBirdID {
			player.Record = updates.Record
			return player, nil
		}
	}
	return nil, errors.New("typie bird does not exist in game room")
}

//Add adds a typie bird into the game room
func (room *GameRoom) Add(bird *TypieBird) error {
	room.Players = append(room.Players, bird)
	return nil
}

//GetByID retrieves the typie bird with `typieBirdID` from the game room
func (room *GameRoom) GetByID(typieBirdID bson.ObjectId) (*TypieBird, error) {
	for _, player := range room.Players {
		if player.ID == typieBirdID {
			return player, nil
		}
	}
	return nil, errors.New("typie bird does not exist in game room")
}

//DeleteByID removes the typie bird with `typieBirdID` from the game room
func (room *GameRoom) DeleteByID(typieBirdID bson.ObjectId) error {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].ID == typieBirdID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			return nil
		}
	}
	return errors.New("typie bird does not exist in game room")
}
