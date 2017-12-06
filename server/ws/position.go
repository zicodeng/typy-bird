// package ws

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/info344-a17/typy-bird/server/handlers"
// 	"gopkg.in/mgo.v2/bson"
// )

// //PositionHandler handles requests for the /position resource
// type PositionHandler struct {
// 	notifier *Notifier
// 	context  *handlers.HandlerContext
// }

// //NewPositionHandler constructs a new PositionHandler
// func NewPositionHandler(notifier *Notifier, context *handlers.HandlerContext) *PositionHandler {
// 	return &PositionHandler{notifier, context}
// }

// //ServeHTTP handles HTTP requests for the UpdateHandler
// func (ph *PositionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	//get ID of current typie bird
// 	queryParams := r.URL.Query()
// 	typieBirdID := bson.ObjectId(queryParams.Get("auth"))

// 	//check current bird is a player (authorize)
// 	if _, err := ph.context.GameRoom.GetByID(typieBirdID); err != nil {
// 		http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	//update position of typie bird
// 	if _, err := ph.context.GameRoom.IncrementPosition(typieBirdID); err != nil {
// 		http.Error(w, fmt.Sprintf("error updating typie bird position: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	//marshall GameRoom struct into json
// 	room, err := json.Marshal(ph.context.GameRoom)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error marshalling gameroom to JSON: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	//respond to client
// 	ph.notifier.Notify(room)
// }
