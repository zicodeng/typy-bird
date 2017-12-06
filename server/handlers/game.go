package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/info344-a17/typy-bird/server/models"
	"gopkg.in/mgo.v2/bson"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	GameRoom   *models.GameRoom
	TypieStore *models.MongoStore
}

//NewHandlerContext creates a new instance of a context struct to be used by a handler
func NewHandlerContext(gameRoom *models.GameRoom, typieStore *models.MongoStore) *HandlerContext {
	return &HandlerContext{
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
		err := json.NewDecoder(r.Body).Decode(newTypie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		//convert new typie bird to typie bird
		typie := newTypie.ToTypie()

		//insert typie bird into the mongo store
		if _, err := c.TypieStore.InsertTypieBird(typie); err != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", err), http.StatusInternalServerError)
			return
		}

		//add typie bird to gameroom
		c.GameRoom.Add(typie)

		//respond to client with created typie bird
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(typie)
		if err != nil {
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
	//get ID from auth header
	queryParams := r.URL.Query()
	typieBirdID := bson.ObjectId(queryParams.Get("auth"))

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

func (c *HandlerContext) DictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		randDictionary := GetRandomDict()
		w.Header().Add("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(randDictionary)
		if err != nil {
			http.Error(w, "error encoding dictionary: %v", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func GetRandomDict() []string {
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
	for i, _ := range perm {
		randDictionary = append(randDictionary, dictionary[i])
	}
	return randDictionary
}
