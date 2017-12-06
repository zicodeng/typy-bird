package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/info344-a17/typy-bird/server/handlers"
	"gopkg.in/mgo.v2/bson"
)

//PositionHandler handles requests for the /position resource
type PositionHandler struct {
	notifier *Notifier
	context  *handlers.HandlerContext
}

//NewPositionHandler constructs a new PositionHandler
func NewPositionHandler(notifier *Notifier, context *handlers.HandlerContext) *PositionHandler {
	return &PositionHandler{notifier, context}
}

//ServeHTTP handles HTTP requests for the UpdateHandler
func (ph *PositionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check current bird is player (authorize)
	queryParams := r.URL.Query()
	typieBirdID := bson.ObjectId(queryParams.Get("auth"))
	if _, err := ph.context.GameRoom.GetByID(typieBirdID); err != nil {
		http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusInternalServerError)
		return
	}

	//marshall GameRoom struct into json
	room, err := json.Marshal(ph.context.GameRoom)
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshalling gameroom to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	//respond to client
	ph.notifier.Notify(room)
}
