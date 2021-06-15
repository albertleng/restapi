package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/albertleng/restapi/file"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	Router *mux.Router
}

var isStarted bool

const (
	baseURL = "https://gateway.marvel.com/v1/public"
	limit   = "100"
)

type config struct {
	privateKey string
	publicKey  string
}

var conf *config

type CharacterId struct {
	Data struct {
		Results []struct {
			Id int `json:"id"`
		} `json:"results"`
	} `json:"data"`
}

type Character struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Results []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Desc string `json:"description"`
		} `json:"results"`
	} `json:"data"`
}

type CodeStatus struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

// Serve an endpoint /characters that returns all the Marvel character ids in a JSON array of
// numbers.
func (a *App) getCharacters(w http.ResponseWriter, r *http.Request) {
	if ids, err := file.ReadFile(); ids != nil && err == nil {
		if isStarted == false {
			isStarted = true
			idsLen := len(ids)
			newIds := getUpdatedCharacterIds(idsLen)
			if len(newIds) > 0 {
				for _, newId := range newIds {
					ids = append(ids, newId)
				}
				file.AppendFile(newIds)
			}
		}

		output, err := json.Marshal(ids)
		if err != nil {
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(output))
		return
	}

	data := getUpdatedCharacterIds(0)

	file.CreateFile()
	file.WriteFile(data)
	output, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
}

func getUpdatedCharacterIds(offset int) (ids []int) {
	var _ []int
	for {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		hash := getMd5(ts + conf.privateKey + conf.publicKey)

		response, err := http.Get(baseURL + "/characters?ts=" + ts + "&apikey=" + conf.publicKey + "&hash=" + hash + "&limit=" + limit + "&offset=" + strconv.Itoa(offset))
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var characterId CharacterId
		err = json.Unmarshal(responseBytes, &characterId)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, result := range characterId.Data.Results {
			ids = append(ids, result.Id)
		}

		if characterId.Data.Results == nil || len(characterId.Data.Results) < 100 {
			break
		}

		offset += 100
	}

	return ids
}

// Serve an endpoint /characters/{characterId} that returns only the id, name and description
// of the character.
func (a *App) getCharacter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets params

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	hash := getMd5(ts + conf.privateKey + conf.publicKey)

	response, err := http.Get(baseURL + "/characters/" + params["characterId"] + "?ts=" + ts + "&apikey=" + conf.publicKey + "&hash=" + hash)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var character Character
	err = json.Unmarshal(responseBytes, &character)
	if err != nil {
		log.Fatal(err)
		return
	}

	if character.Status != "Ok" {
		codeStatus := CodeStatus{Status: character.Status, Code: character.Code}
		data, _ := json.MarshalIndent(codeStatus, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(data))
		return
	}

	data, err := json.MarshalIndent(character.Data.Results[0], "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/characters", a.getCharacters).Methods("GET")
	a.Router.HandleFunc("/characters/{characterId}", a.getCharacter).Methods("GET")
}

func readConfig() *config {
	return &config{
		privateKey: os.Getenv("MARVEL_API_PRIVATE_KEY"),
		publicKey:  os.Getenv("MARVEL_API_PUBLIC_KEY"),
	}
}

func getMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
	})
	handler := c.Handler(a.Router)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func (a *App) Initialize() {
	conf = readConfig()

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}
