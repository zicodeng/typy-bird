package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/info344-a17/typy-bird/server/models"
	"gopkg.in/mgo.v2/bson"
)

//GameroomHandler handles the /gameroom route and returns the current gameroom
func (c *HandlerContext) GameroomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(c.GameRoom); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		//get ID of current typie bird
		queryParams := r.URL.Query()
		typieBirdID := bson.ObjectIdHex(queryParams.Get("auth"))

		//check current bird is a player in the game room (authorize)
		if _, err := c.GameRoom.GetByID(typieBirdID); err != nil {
			http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusBadRequest)
			return
		}

		err := c.GameRoom.Delete(typieBirdID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error deleting bird from players: %v", err), http.StatusInternalServerError)
			return
		}

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"Position",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//StartGameHandler starts a new game
func (c *HandlerContext) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.GameRoom.Available = false

		startTime := time.Now()
		c.GameRoom.StartTime = startTime

		wsPayload := struct {
			Type      string    `json:"type,omitempty"`
			StartTime time.Time `json:"startTime,omitempty"`
		}{
			"GameStart",
			startTime,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//EndGameHandler ends the current game
func (c *HandlerContext) EndGameHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//reset gameroom to initial state
		c.GameRoom.Available = true
		for _, client := range c.GameRoom.Players {
			c.GameRoom.Delete(client.ID)
		}

		//broadcast new gameroom state to client
		gwsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"GameEnd",
			c.GameRoom,
		}
		gamePayload, jsonErr := json.Marshal(gwsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(gamePayload)

		//get top scorers from mongo store
		leaders, getErr := c.TypieStore.GetTopScores()
		if getErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", getErr), http.StatusInternalServerError)
			return
		}

		//create LeaderBoard struct and marshall to json
		leaderBoard := &models.LeaderBoard{
			Leaders:   leaders,
			Available: c.GameRoom.Available,
		}

		//broadcast new gameroom state to client
		lwsPayload := struct {
			Type    string              `json:"type,omitempty"`
			Payload *models.LeaderBoard `json:"payload,omitempty"`
		}{
			"Leaderboard",
			leaderBoard,
		}
		leaderPayload, jsonErr := json.Marshal(lwsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(leaderPayload)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}
