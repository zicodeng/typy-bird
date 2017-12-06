package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/info344-a17/typy-bird/server/handlers"
	"github.com/info344-a17/typy-bird/server/models"
)

//LeaderBoardHandler handles requests for the /position resource
type LeaderBoardHandler struct {
	notifier *Notifier
	context  *handlers.HandlerContext
}

//NewLeaderBoardHandler constructs a new PositionHandler
func NewLeaderBoardHandler(notifier *Notifier, context *handlers.HandlerContext) *LeaderBoardHandler {
	return &LeaderBoardHandler{notifier, context}
}

//ServeHTTP handles HTTP requests for the UpdateHandler
func (lh *LeaderBoardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//get top scorers from mongo store
	leaders, getErr := lh.context.TypieStore.GetTopScores()
	if getErr != nil {
		http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", getErr), http.StatusInternalServerError)
		return
	}

	//create LeaderBoard struct and marshall into json
	leaderBoard := &models.LeaderBoard{
		Leaders:   leaders,
		Available: lh.context.GameRoom.Available,
	}
	board, jsonErr := json.Marshal(leaderBoard)
	if jsonErr != nil {
		http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", jsonErr), http.StatusInternalServerError)
		return
	}

	//respond to client
	lh.notifier.Notify(board)
}
