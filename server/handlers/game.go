package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/info344-a17/typy-bird/server/models"
	"gopkg.in/mgo.v2/bson"
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

//TypieHandler handles methods for the /typie route
func (c *HandlerContext) TypieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//decode new typie bird from request body
		newTypie := &models.NewTypieBird{}
		if err := json.NewDecoder(r.Body).Decode(newTypie); err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		//convert new typie bird to typie bird
		typie := newTypie.ToTypie()

		//insert typie bird into the mongo store
		insertedTypie, insErr := c.TypieStore.InsertTypieBird(typie)
		if insErr != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", insErr), http.StatusInternalServerError)
			return
		}

		//add typie bird to gameroom
		c.GameRoom.Add(insertedTypie)

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"NewTypie",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with created typie bird
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(insertedTypie); err != nil {
			http.Error(w, fmt.Sprintf("error encoding the created typie: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//TypieMeHandler handles the methods for the /typie/me route
func (c *HandlerContext) TypieMeHandler(w http.ResponseWriter, r *http.Request) {
	//get bird associated with current ID
	queryParams := r.URL.Query()
	typieBirdID := bson.ObjectIdHex(queryParams.Get("auth"))

	switch r.Method {
	case "GET":
		//get bird associate with current ID
		bird, err := c.TypieStore.GetByID(typieBirdID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error retrieving typie bird from store: %v", err), http.StatusBadRequest)
			return
		}

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "PATCH":
		//decode new record from request body
		updates := &models.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in gameroom
		bird, err := c.GameRoom.Update(typieBirdID, updates)
		if err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in typie store
		if err := c.TypieStore.Update(typieBirdID, updates); err != nil {
			http.Error(w, fmt.Sprintf("error updating user store: %v", err), http.StatusBadRequest)
			return
		}

		//get top scorers from mongo store
		leaders, getErr := c.TypieStore.GetTopScores()
		if getErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", getErr), http.StatusInternalServerError)
			return
		}

		//create LeaderBoard struct and marshall to json
		leaderBoard := &models.LeaderBoard{
			Leaders:   leaders,
			Available: c.GameRoom.Available,
		}

		wsPayload := struct {
			Type    string              `json:"type,omitempty"`
			Payload *models.LeaderBoard `json:"payload,omitempty"`
		}{
			"Leaderboard",
			leaderBoard,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		//remove typie bird from game room
		c.GameRoom.Delete(typieBirdID)

		//respond to client
		w.Write([]byte("game ended for player\n"))
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//PositionHandler handles the /position route and returns the current postion of a bird
func (c *HandlerContext) PositionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		//get ID of current typie bird
		queryParams := r.URL.Query()
		typieBirdID := bson.ObjectId(queryParams.Get("auth"))

		//check current bird is a player in the game room (authorize)
		if _, err := c.GameRoom.GetByID(typieBirdID); err != nil {
			http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusBadRequest)
			return
		}

		//update position of typie bird in gameroom struct
		bird, incrErr := c.GameRoom.IncrementPosition(typieBirdID)
		if incrErr != nil {
			http.Error(w, fmt.Sprintf("error updating typie bird position: %v", incrErr), http.StatusInternalServerError)
			return
		}

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"Position",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//GameroomHandler handles the /gameroom route and returns the current gameroom
func (c *HandlerContext) GameroomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(c.GameRoom); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	case "POST":
		startTime := time.Now()

		wsPayload := struct {
			Type      string    `json:"type,omitempty"`
			StartTime time.Time `json:"startTime,omitempty"`
		}{
			"GameStart",
			startTime,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//ReadyHandler handles updating the typies ready status
func (c *HandlerContext) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		//get ID of current typie bird
		queryParams := r.URL.Query()
		typieBirdID := bson.ObjectId(queryParams.Get("auth"))

		//check current bird is a player in the game room (authorize)
		if _, err := c.GameRoom.GetByID(typieBirdID); err != nil {
			http.Error(w, fmt.Sprintf("error getting typie bird: %v", err), http.StatusBadRequest)
			return
		}

		//update status of typie bird in gameroom struct
		bird, incrErr := c.GameRoom.ReadyUp(typieBirdID)
		if incrErr != nil {
			http.Error(w, fmt.Sprintf("error updating typie bird ready status: %v", incrErr), http.StatusInternalServerError)
			return
		}

		wsPayload := struct {
			Type    string           `json:"type,omitempty"`
			Payload *models.GameRoom `json:"payload,omitempty"`
		}{
			"Ready",
			c.GameRoom,
		}
		//broadcast new gameroom state to client
		payload, jsonErr := json.Marshal(wsPayload)
		if jsonErr != nil {
			http.Error(w, fmt.Sprintf("error marshalling payload to JSON: %v", jsonErr), http.StatusInternalServerError)
			return
		}
		c.Notifier.Notify(payload)

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//LeaderboardHandler retrieves the top 10 in ascending order
func (c *HandlerContext) LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//get top scorers from mongo store
		leaderboard, err := c.TypieStore.GetTopScores()
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling leaderboard to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(leaderboard); err != nil {
			http.Error(w, fmt.Sprintf("error encoding leaderboard to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

//DictHandler handles the /dictionary route and gets a list of random words that the users must type to complete the game
func (c *HandlerContext) DictHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		randDictionary := getRandomDict()
		w.Header().Add("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(randDictionary)
		if err != nil {
			http.Error(w, "error encoding dictionary: %v", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func getRandomDict() []string {
	fourLetterWords := [25]string{"curl", "etas", "pleb", "tabi", "soup", "tune", "kure", "tech",
		"suez", "veld", "bash", "cole", "peek", "kill", "tarn", "momi", "flee", "cone",
		"cham", "land", "amok", "ship", "maim", "bird", "prig"}
	fiveLetterWords := [25]string{"roost", "gayly", "ptain", "unbid", "umiac", "kappa", "festa", "every",
		"playa", "olden", "donna", "godin", "muzio", "sauce", "blink", "cause", "chirm", "oriel", "schwa",
		"bogle", "chick", "dolin", "loads", "using", "sweal"}
	sixLetterWords := [25]string{"catton", "khedah", "untold", "bhindi", "decree", "kinase", "cohere",
		"waffie", "garter", "bashan", "roddie", "stingo", "dodger", "chalet", "contra", "blanch",
		"edwina", "immesh", "fulmar", "saddle", "finish", "piggin", "riches", "dengue", "mizzle"}
	sevenLetterWords := [25]string{"unmined", "rosario", "ericoid", "herbert", "faraway", "grimace",
		"brioche", "napless", "deprive", "inhered", "plantin", "outpour", "whoosis", "impanel",
		"stuffed", "taussig", "narvez", "seattle", "millier", "leister", "arduous", "ransome",
		"luzerne", "bunches", "bighead"}
	dictArray := [4][25]string{fourLetterWords, fiveLetterWords, sixLetterWords, sevenLetterWords}
	dictionary := []string{}
	randDictionary := []string{}
	perm := rand.Perm(25)
	for i := 0; i < len(dictArray); i++ {
		for index, value := range perm {
			if index == 5 {
				break
			}
			dictionary = append(dictionary, dictArray[i][value])
		}
	}
	perm = rand.Perm(20)
	for _, v := range perm {
		randDictionary = append(randDictionary, dictionary[v])
	}
	return randDictionary
}
