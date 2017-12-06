package main

import (
	"log"
	"net/http"
	"os"

	"github.com/info344-a17/typy-bird/server/handlers"
	"github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/ws"
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
	gameRoom := &models.GameRoom{Available: true}

	//Initialize handler stuff
	context := handlers.NewHandlerContext(gameRoom, typieStore)
	notifier := ws.NewNotifier()
	mux := http.NewServeMux()

	//POST for creating new typies
	mux.HandleFunc("/typie", context.TypieHandler)
	//GET for dictionary words
	mux.HandleFunc("/dictionary", context.DictHandler)
	//PATCH for updating typie records
	mux.HandleFunc("/typie/me", context.TypieMeHandler)
	//GET for gameroom
	mux.HandleFunc("/gameroom", context.GameroomHandler)
	//upgrading to websockets
	mux.Handle("/ws", ws.NewWebSocketsHandler(notifier))
	//sending postions to players
	mux.Handle("/position", ws.NewPositionHandler(notifier, context))
	//sending leader stats to waiting players
	mux.Handle("/leaderboard", ws.NewLeaderBoardHandler(notifier, context))

	corsMux := handlers.NewCORSHandler(mux)

	log.Printf("server is listening at http://%s...", addr)
	log.Fatal(http.ListenAndServe(addr, corsMux))
}
