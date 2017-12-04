package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ezhai24/challenges-ezhai24/servers/gateway/sessions"
)

//TypieMeHandler handlers current birds
func (ctx *HandlerContext) TypieMeHandler(w http.ResponseWriter, r *http.Request) {
	//get session state associated with current typie bird
	sessID, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &SessionState{})
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "PATCH":
		//decode new record from request body
		record := &game.Record{}
		if err := json.NewDecoder(r.Body).Decode(record); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in state
		if err := state.TypieBird.UpdateRecord(record); err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}

		// update bird in session store
		if err := ctx.SessionStore.Save(sessID, state); err != nil {
			http.Error(w, fmt.Sprintf("error saving to session store: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in typie store
		if err := ctx.TypieStore.Update(state.TypieBird.UserName, record); err != nil {
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
