package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/sessions"
)

//TypieMeHandler handlers current birds
func (ctx *HandlerContext) TypieMeHandler(w http.ResponseWriter, r *http.Request) {
	//get session state associated with current typie bird
	state := &SessionState{}
	sessID, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "PATCH":
		//decode new record from request body
		updates := &game.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in state
		if err := state.TypieBird.UpdateRecord(updates); err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}

		// update bird in session store
		if err := ctx.SessionStore.Save(sessID, state); err != nil {
			http.Error(w, fmt.Sprintf("error saving to session store: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in typie store
		if err := ctx.TypieStore.Update(state.TypieBird.ID, updates); err != nil {
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
