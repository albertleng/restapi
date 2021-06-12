package main

// Ref: https://github.com/loivis/marvel-comics-api-data-loader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func getCharacters(w http.ResponseWriter, r *http.Request) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	hash := getMd5(ts + conf.privateKey + conf.publicKey)

	response, err := http.Get(baseURL + "/characters?ts=" + ts + "&apikey=" + conf.publicKey + "&hash=" + hash + "&limit=" + limit)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fmt.Println("ts: " + ts)
	fmt.Println("apikey: " + conf.publicKey)
	fmt.Println("hash: " + hash)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, responseString)
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

	fmt.Println("ts: " + ts)
	fmt.Println("apikey: " + conf.publicKey)
	fmt.Println("hash: " + hash)
	fmt.Println("characterId: " + params["characterId"])

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, responseString)

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
