package handlers

import (
	"math/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/info344-a17/typy-bird/server/models"
)

//HandlerContext keeps track of database information
type HandlerContext struct {
	SessionKey   string
	TypieStore   *models.MongoStore
}

//NewHandlerContext creates a new instance of a context struct to be used by a handler
func NewHandlerContext(key string, typieStore *models.MongoStore) *HandlerContext {
	return &HandlerContext{
		SessionKey:   key,
		TypieStore:   typieStore,
	}
}

//SessionState keeps track of current session information
type SessionState struct {
	SessionStart time.Time
	TypieBird    *models.TypieBird
}

//TypieHandler handles the POST,GET, and PATCH methods for the /typie route
func (c *HandlerContext) TypieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		newTypie := &models.NewTypieBird{}
		err := json.NewDecoder(r.Body).Decode(newTypie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding typie json: %v", err), http.StatusInternalServerError)
			return
		}

		typie := newTypie.ToTypie()

		insertedTypie, err := c.TypieStore.InsertTypieBird(typie)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting typie: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(insertedTypie)
		if err != nil {
			http.Error(w, "error encoding the created typie", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case "GET":
		leaderboard, err := c.TypieStore.GetTopScores()
		if err != nil {
			http.Error(w, fmt.Sprintf("error retrieving leaderboard: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		err = json.NewEncoder(w).Encode(leaderboard)
		if err != nil {
			http.Error(w, "error encoding leaderboard: %v", http.StatusInternalServerError)
			return
		}
	case "PATCH":
		//get session state associated with current typie bird
		state := &SessionState{}

		//decode new record from request body
		updates := &models.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		//update bird in state
		if err := state.TypieBird.UpdateRecord(updates); err != nil {
			http.Error(w, fmt.Sprintf("error applying updates: %v", err), http.StatusBadRequest)
			return
		}


		//update bird in typie store
		if err := c.TypieStore.Update(state.TypieBird.ID, updates); err != nil {
			http.Error(w, fmt.Sprintf("error updating user store: %v", err), http.StatusBadRequest)
			return
		}

		//respond to client with updated bird
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(state.TypieBird); err != nil {
			http.Error(w, fmt.Sprintf("error encoding user to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func (c *HandlerContext) DictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fourLetterWords := [25]string{"curl", "etas","pleb","tabi","soup","tune","kure","tech",
			"suez","veld","bash","cole","peek","kill","tarn","momi","flee","cone",
			"cham","land","amok","ship","maim","bird","prig"}
		fiveLetterWords := [25]string{"roost","gayly","ptain","unbid","umiac","kappa","festa","every",
			"playa","olden","donna","godin","muzio","sauce","blink","cause","chirm","oriel","schwa",
			"bogle","chick","dolin","loads","using","sweal"}
		sixLetterWords := [25]string{"catton","khedah","untold","bhindi","decree","kinase","cohere",
			"waffie","garter","bashan","roddie","stingo","dodger","chalet","contra","blanch",
			"edwina","immesh","fulmar","saddle","finish","piggin","riches","dengue","mizzle",}
		sevenLetterWords := [25]string{"unmined","rosario","ericoid","herbert","faraway","grimace",
			"brioche","napless","deprive","inhered","plantin","outpour","whoosis","impanel",
			"stuffed","taussig","narvez","seattle","millier","leister","arduous","ransome",
			"luzerne","bunches","bighead"}
		dictArray := [4][25]string{fourLetterWords, fiveLetterWords, sixLetterWords, sevenLetterWords}
		dictionary := make([]string, 20)
		randDictionary := make([]string, len(dictionary))
		perm := rand.Perm(25)
		count := 0
		for i := 0; i < 4; i++ {
			for index, value := range perm {
				if index == 5 {
					break
				}
				dictionary[count] = dictArray[i][value]
				count++
			}
		}
		perm = rand.Perm(20)
		for i, v := range perm {
			randDictionary[v] = dictionary[i]
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

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
