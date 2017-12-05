package ws

import (
	"net/http"
	"strconv"

	"github.com/info344-a17/typy-bird/server/handlers"
	"github.com/info344-a17/typy-bird/server/sessions"
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
	state := &handlers.SessionState{}
	sessions.GetState(r, ph.context.SessionKey, ph.context.SessionStore, state)
	ph.notifier.Notify([]byte(strconv.Itoa(state.TypieBird.Position)))
}
