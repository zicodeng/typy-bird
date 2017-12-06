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

	//Initialize handler stuff
	context := handlers.NewHandlerContext(typieStore)
	notifier := ws.NewNotifier()
	mux := http.NewServeMux()

	//POST,GET,PATCH for typies
	mux.HandleFunc("/typie", context.TypieHandler)
	//GET for dictionary words
	mux.HandleFunc("/dictionary", context.DictHandler)
	//upgrading to websockets
	mux.Handle("/ws", ws.NewWebSocketsHandler(notifier))
	//sending postions to players
	mux.Handle("/update", ws.NewUpdateHandler(notifier))

	log.Printf("server is listening at http://%s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
