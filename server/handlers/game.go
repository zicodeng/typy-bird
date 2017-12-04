package handlers

import (
	"time"

	"github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/sessions"
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
