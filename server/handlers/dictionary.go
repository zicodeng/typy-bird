package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

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
