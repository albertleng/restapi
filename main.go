package main

// Ref: https://github.com/loivis/marvel-comics-api-data-loader

//TODO: Implement unit testing
//TODO: Write a swagger.yaml
//TODO: Add readme.md

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/albertleng/restapi/file"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	baseURL = "https://gateway.marvel.com/v1/public"
	limit   = "100"
)

type CharacterId struct {
	Data struct {
		Results []struct {
			Id int `json:"id"`
		} `json:"results"`
	} `json:"data"`
}

type Character struct {
	Data struct {
		Results []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Desc string `json:"description"`
		} `json:"results"`
	} `json:"data"`
}

// TODO: 1. Get ALL characters by multiple calls using different offsets
// TODO: 2. Implement caching
// Serve an endpoint /characters that returns all the Marvel character ids in a JSON array of
// numbers.
func getCharacters(w http.ResponseWriter, _ *http.Request) {

	if ids, err := file.ReadFile(); ids != nil && err == nil {
		output, err := json.Marshal(ids)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Fprint(w, string(output))

		return
	}

	var data []int
	offset := 0
	for {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		hash := getMd5(ts + conf.privateKey + conf.publicKey)

		response, err := http.Get(baseURL + "/characters?ts=" + ts + "&apikey=" + conf.publicKey + "&hash=" + hash + "&limit=" + limit + "&offset=" + strconv.Itoa(offset))
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		fmt.Println("ts: " + ts)
		fmt.Println("apikey: " + conf.publicKey)
		fmt.Println("hash: " + hash)
		fmt.Println("offset: " + strconv.Itoa(offset))

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
			data = append(data, result.Id)
		}

		fmt.Println("characterId.Data.Results == nil: " + strconv.FormatBool(characterId.Data.Results == nil))
		fmt.Println("characterId.Data.Results < 100: " + strconv.FormatBool(len(characterId.Data.Results) < 100))
		fmt.Println("len(characterId.Data.Results): " + strconv.FormatInt(int64(len(characterId.Data.Results)), 10))
		if characterId.Data.Results == nil || len(characterId.Data.Results) < 100 {
			break
		}

		offset += 100
	}
	file.CreateFile()
	file.WriteFile(data)
	output, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprint(w, string(output))
}

// Serve an endpoint /characters/{characterId} that returns only the id, name and description
// of the character.
func getCharacter(w http.ResponseWriter, r *http.Request) {
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

	data, err := json.MarshalIndent(character.Data.Results[0], "", " ")
	fmt.Fprint(w, string(data))
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/characters", getCharacters).Methods("GET")
	router.HandleFunc("/characters/{characterId}", getCharacter).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type config struct {
	privateKey string
	publicKey  string
}

// TODO: Change to read from config file
func readConfig() *config {
	return &config{
		privateKey: os.Getenv("MARVEL_API_PRIVATE_KEY"),
		publicKey:  os.Getenv("MARVEL_API_PUBLIC_KEY"),
	}
}

var conf *config

func getMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	conf = readConfig()
	handleRequests()
}
