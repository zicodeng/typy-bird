package handlers

import (
	"github.com/info344-a17/typy-bird/server/models"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	Notifier   *Notifier
	GameRoom   *models.GameRoom
	TypieStore *models.MongoStore
}

//NewHandlerContext creates a new instance of a context struct to be used by a handler
func NewHandlerContext(notifier *Notifier, gameRoom *models.GameRoom, typieStore *models.MongoStore) *HandlerContext {
	return &HandlerContext{
		Notifier:   notifier,
		GameRoom:   gameRoom,
		TypieStore: typieStore,
	}
}
