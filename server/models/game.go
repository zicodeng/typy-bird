package models

import (
	"errors"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//ErrTypieBirdNotPlayer is thrown when a typie bird with the given ID is not in the game room player struct
var ErrTypieBirdNotPlayer = errors.New("typie bird does not exist in game room")

//LeaderBoard represents highest scoring typie birds
type LeaderBoard struct {
	Leaders   []*TypieBird `json:"leaders,omitempty"`
	Available bool         `json:"available,omitempty"`
}

//GameRoom represents the room of players
type GameRoom struct {
	Players   []*TypieBird `json:"players,omitempty"`
	Available bool         `json:"available,omitempty"`
	StartTime time.Time    `json:"startTime,omitempty"`
	mx        sync.RWMutex
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
	IsReady  bool          `json:"isReady"`
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
	room.mx.Lock()
	defer room.mx.Unlock()

	for _, player := range room.Players {
		if player.ID == typieBirdID {
			player.Record = updates.Record
			return player, nil
		}
	}
	return nil, ErrTypieBirdNotPlayer
}

//Add adds a typie bird into the game room
func (room *GameRoom) Add(bird *TypieBird) error {
	//check room is not already full
	if len(room.Players) > 4 {
		return errors.New("gameroom full")
	}

	//protect with mutex
	room.mx.Lock()
	defer room.mx.Unlock()

	//add typie bird to game room
	room.Players = append(room.Players, bird)

	//change gameroom availability if necessary
	if len(room.Players) == 4 {
		room.Available = false
	}

	return nil
}

//GetByID retrieves the typie bird with `typieBirdID` from the game room
func (room *GameRoom) GetByID(typieBirdID bson.ObjectId) (*TypieBird, error) {
	room.mx.RLock()
	defer room.mx.RUnlock()

	for _, player := range room.Players {
		if player.ID == typieBirdID {
			return player, nil
		}
	}
	return nil, ErrTypieBirdNotPlayer
}

//Delete removes the typie bird with `typieBirdID` from the game room
func (room *GameRoom) Delete(typieBirdID bson.ObjectId) error {
	room.mx.Lock()
	defer room.mx.Unlock()

	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].ID == typieBirdID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			return nil
		}
	}
	return ErrTypieBirdNotPlayer
}

//IncrementPosition increments the position of the given bird by one step
func (room *GameRoom) IncrementPosition(typieBirdID bson.ObjectId) (*TypieBird, error) {
	room.mx.Lock()
	defer room.mx.Unlock()

	for _, player := range room.Players {
		if player.ID == typieBirdID {
			player.Position = player.Position + 1
			return player, nil
		}
	}
	return nil, ErrTypieBirdNotPlayer
}

//ReadyUp changes a bird's status from not ready to ready
func (room *GameRoom) ReadyUp(typieBirdID bson.ObjectId) (*TypieBird, error) {
	room.mx.Lock()
	defer room.mx.Unlock()

	for _, player := range room.Players {
		if player.ID == typieBirdID {
			player.IsReady = !player.IsReady
			return player, nil
		}
	}
	return nil, ErrTypieBirdNotPlayer
}
