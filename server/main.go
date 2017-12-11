package main

import (
	"log"
	"net/http"
	"os"

	"github.com/info344-a17/typy-bird/server/handlers"
	"github.com/info344-a17/typy-bird/server/models"
	"gopkg.in/mgo.v2"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	//MongoDB Initialization
	dbAddr := os.Getenv("DBADDR")
	if len(dbAddr) == 0 {
		log.Fatal("the DBADDR was not set")
	}
	mongoSess, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("error dialing db: %v", err)
	}
	typieStore := models.NewMongoStore(mongoSess, "GameDB", "TypieCollection")

	//Initialize game room struct
	gameRoom := &models.GameRoom{}
	gameRoom.Available = true

	//Initialize handler stuff
	notifier := handlers.NewNotifier()
	context := handlers.NewHandlerContext(notifier, gameRoom, typieStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/typie", context.TypieHandler)
	mux.HandleFunc("/typie/me", context.TypieMeHandler)
	mux.HandleFunc("/typie/position", context.PositionHandler)

	mux.HandleFunc("/gameroom", context.GameroomHandler)
	mux.HandleFunc("/ready", context.ReadyHandler)
	mux.HandleFunc("/start", context.StartGameHandler)
	mux.HandleFunc("/end", context.EndGameHandler)

	mux.HandleFunc("/dictionary", context.DictHandler)
	mux.HandleFunc("/leaderboard", context.LeaderboardHandler)
	mux.Handle("/ws", handlers.NewWebSocketsHandler(notifier))

	corsMux := handlers.NewCORSHandler(mux)

	log.Printf("server is listening at http://%s...", addr)
	log.Fatal(http.ListenAndServe(addr, corsMux))
}
