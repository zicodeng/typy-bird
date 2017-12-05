package ws

import (
	"net/http"
	"strconv"

	"github.com/info344-a17/typy-bird/server/handlers"
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
	// state := &handlers.SessionState{}
	// if _, err := sessions.GetState(r, ph.context.SessionKey, ph.context.SessionStore, state); err != nil {
	// 	http.Error(w, fmt.Sprintf("error getting state: %v", err), http.StatusInternalServerError)
	// 	return
	// }
	ph.notifier.Notify([]byte(strconv.Itoa(state.TypieBird.Position)))
}
