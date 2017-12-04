package handlers

import (
	"time"

	"github.com/ezhai24/challenges-ezhai24/servers/gateway/sessions"
	"github.com/ezhai24/typy-bird/server/models"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	SessionKey   string
	SessionStore sessions.Store
	TypieStore   *game.MongoStore
}

//SessionState keeps track of current session information
type SessionState struct {
	SessionStart time.Time
	TypieBird    *game.TypieBird
}
