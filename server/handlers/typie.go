package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/info344-a17/typy-bird/server/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

//TypieHandler handles methods for the /typie route
func (c *HandlerContext) TypieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//decode new typie bird from request body
		newTypie := &models.NewTypieBird{}
		if err := json.NewDecoder(r.Body).Decode(newTypie); err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		//convert new typie bird to typie bird
		typie := newTypie.ToTypie()

		//insert typie bird into the mongo store
		insertedTypie, insErr := c.TypieStore.InsertTypieBird(typie)
		if insErr != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", insErr), http.StatusInternalServerError)
			return
		}

		//add typie bird to gameroom
		c.GameRoom.Add(insertedTypie)

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"NewTypie",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with created typie bird
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(insertedTypie); err != nil {
			http.Error(w, fmt.Sprintf("error encoding the created typie: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//TypieMeHandler handles the methods for the /typie/me route
func (c *HandlerContext) TypieMeHandler(w http.ResponseWriter, r *http.Request) {
	//get bird associated with current ID
	queryParams := r.URL.Query()
	typieBirdID := bson.ObjectIdHex(queryParams.Get("auth"))

	switch r.Method {
	case "GET":
		//get bird associate with current ID
		bird, err := c.TypieStore.GetByID(typieBirdID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error retrieving typie bird from store: %v", err), http.StatusBadRequest)
			return
		}

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "PATCH":
		//decode new record from request body
		updates := &models.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in gameroom
		bird, err := c.GameRoom.Update(typieBirdID, updates)
		if err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in typie store
		if err := c.TypieStore.Update(typieBirdID, updates); err != nil {
			http.Error(w, fmt.Sprintf("error updating user store: %v", err), http.StatusBadRequest)
			return
		}

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		//remove typie bird from game room
		c.GameRoom.Delete(typieBirdID)

		//respond to client
		w.Write([]byte("game ended for player\n"))
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//PositionHandler handles the /position route and returns the current postion of a bird
func (c *HandlerContext) PositionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		//get ID of current typie bird
		queryParams := r.URL.Query()
		typieBirdID := bson.ObjectIdHex(queryParams.Get("auth"))

		//check current bird is a player in the game room (authorize)
		if _, err := c.GameRoom.GetByID(typieBirdID); err != nil {
			http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusBadRequest)
			return
		}

		//update position of typie bird in gameroom struct
		bird, incrErr := c.GameRoom.IncrementPosition(typieBirdID)
		if incrErr != nil {
			http.Error(w, fmt.Sprintf("error updating typie bird position: %v", incrErr), http.StatusInternalServerError)
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

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//ReadyHandler handles updating the typies ready status
func (c *HandlerContext) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		//get ID of current typie bird
		queryParams := r.URL.Query()
		typieBirdID := bson.ObjectIdHex(queryParams.Get("auth"))

		//check current bird is a player in the game room (authorize)
		if _, err := c.GameRoom.GetByID(typieBirdID); err != nil {
			http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusBadRequest)
			return
		}

		//update status of typie bird in gameroom struct
		bird, incrErr := c.GameRoom.ReadyUp(typieBirdID)
		if incrErr != nil {
			http.Error(w, fmt.Sprintf("error updating typie bird ready status: %v", incrErr), http.StatusInternalServerError)
			return
		}

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"Ready",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}
