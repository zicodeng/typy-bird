package handlers

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"net/http"

)

//HandlerContext keeps track of database information
type HandlerContext struct {
	SessionStore sessions.Store
	TypieStore   *game.MongoStore
}

//SessionState keeps track of current session information
type SessionState struct {
	SessionStart time.Time
	TypieBird    *game.TypieBird
}

func (c *HandlerContext) TypieHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		newTypie := &models.TypieBird
		err := json.NewDecoder(r.Body).Decode(newTypie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		typie, err := c.TypieStore.InsertTypieBird(newTypie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", err), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(typie)
		if err != nil {
			http.Error(w, "error encoding the created typie", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case "GET":
		r.Header.Add("Content-Type", "application/json")

		leaderboard, err := c.TypieStore.GetAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("error retrieving leaderboard: %v", err), http.StatusInternalServerError)
		}

		err = json.NewEncoder(w).Encode(leaderboard)
		if err != nil {
			http.Error(w, "error encoding leaderboard: %v", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	default:
		http.Error(w, "method must POST or GET", http.StatusMethodNotAllowed)
		return
	}
}