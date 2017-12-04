package models 

import (
	"time"
)


type TypieBird struct {
	UserName  string        `json:"userName"`
	Record time.Time
}

type Updates struct {
	Record	time.Time
}
