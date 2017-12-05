package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/sessions"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	SessionKey   string
	SessionStore sessions.Store
	TypieStore   *models.MongoStore
}

//NewHandlerContext creates a new instance of a context struct to be used by a handler
func NewHandlerContext(key string, sessionStore sessions.Store, typieStore *models.MongoStore) *HandlerContext {
	return &HandlerContext{
		SessionKey:   key,
		SessionStore: sessionStore,
		TypieStore:   typieStore,
	}
}

//SessionState keeps track of current session information
type SessionState struct {
	SessionStart time.Time
	TypieBird    *models.TypieBird
}

//TypieHandler handles methods for the /typie route
func (c *HandlerContext) TypieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//decode new typie bird from request body
		newTypie := &models.NewTypieBird{}
		err := json.NewDecoder(r.Body).Decode(newTypie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		//check username does not exist in mongostore
		if _, err := c.TypieStore.GetByUserName(newTypie.UserName); err == nil {
			http.Error(w, fmt.Sprintf("a typie bird with this username already exists: %v", err), http.StatusBadRequest)
			return
		}

		//convert new typie bird to typie bird
		typie := newTypie.ToTypie()

		//insert typie bird into the mongo store
		insertedTypie, err := c.TypieStore.InsertTypieBird(typie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", err), http.StatusInternalServerError)
			return
		}

		//respond to client
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(insertedTypie)
		if err != nil {
			http.Error(w, "error encoding the created typie", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case "GET":
		r.Header.Add("Content-Type", "application/json")

		leaderboard, err := c.TypieStore.GetTopScores()
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
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//TypieMeHandler handles the methods for the /typie/me route
func (c *HandlerContext) TypieMeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		//get session state associated with current typie bird
		state := &SessionState{}
		sessID, err := sessions.GetState(r, c.SessionKey, c.SessionStore, state)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusUnauthorized)
			return
		}

		//decode new record from request body
		updates := &models.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in state
		if err := state.TypieBird.Update(updates); err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}

		// update bird in session store
		if err := c.SessionStore.Save(sessID, state); err != nil {
			http.Error(w, fmt.Sprintf("error saving to session store: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in typie store
		if err := c.TypieStore.Update(state.TypieBird.ID, updates); err != nil {
			http.Error(w, fmt.Sprintf("error updating user store: %v", err), http.StatusBadRequest)
			return
		}

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(state.TypieBird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handles the methods for the /sessions route
func (c *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//decode credentials from request body
		credentials := &models.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(credentials); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//check username exists in store (authenticate bird)
		bird, err := c.TypieStore.GetByUserName(credentials.UserName)
		if err != nil {
			http.Error(w, fmt.Sprintf("a typie bird with this username does not exist: %v", err), http.StatusBadRequest)
			return
		}

		//begin a new session
		state := &SessionState{
			SessionStart: time.Now(),
			TypieBird:    bird,
		}
		if _, err := sessions.BeginSession(c.SessionKey, c.SessionStore, state, w); err != nil {
			http.Error(w, fmt.Sprintf("error starting session: %v", err), http.StatusInternalServerError)
			return
		}

		//respond to client with the session bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding typie bird to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "GET":
		//get all current running sessions
		sessions, err := c.SessionStore.GetAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("error retrieving current sessions: %v", err), http.StatusInternalServerError)
		}

		//respond to client
		r.Header.Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sessions); err != nil {
			http.Error(w, "error encoding leaderboard: %v", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}
