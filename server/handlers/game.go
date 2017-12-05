package handlers

import (
	"time"
  "gopkg.in/mgo.v2/bson"
	"fmt"
	"net/http"
  
  "github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/sessions"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	SessionKey   string
	SessionStore sessions.Store
	TypieStore   *models.MongoStore
}

//SessionState keeps track of current session information
type SessionState struct {
	SessionStart time.Time
	TypieBird    *models.TypieBird
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