package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/info344-a17/typy-bird/server/handlers"
	"github.com/info344-a17/typy-bird/server/models"
	"github.com/info344-a17/typy-bird/server/sessions"
	"github.com/info344-a17/typy-bird/server/ws"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	//Session Initialization
	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("the SESSIONKEY was not set")
	}
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		log.Fatal("the REDISADDR was not set")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	redisStore := sessions.NewRedisStore(redisClient, time.Minute*15)

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

	//POST,GET,PATCH,DELETE for typies
	mux.HandleFunc("/typie", context.TypieHandler)
	mux.HandleFunc("/typie/me", context.TypieMeHandler)
	//upgrading to websockets
	mux.Handle("/ws", ws.NewWebSocketsHandler(notifier))
	//sending postions to players
	mux.Handle("/position", ws.NewPositionHandler(notifier, context))

	log.Printf("server is listening at http://%s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
