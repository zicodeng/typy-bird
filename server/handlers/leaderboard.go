package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//LeaderboardHandler retrieves the top 10 in ascending order
func (c *HandlerContext) LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//get top scorers from mongo store
		leaderboard, err := c.TypieStore.GetTopScores()
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(leaderboard); err != nil {
			http.Error(w, fmt.Sprintf("error encoding leaderboard to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}
