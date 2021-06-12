package main

// Ref: https://github.com/loivis/marvel-comics-api-data-loader

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	baseURL = "https://gateway.marvel.com/v1/public"
	limit = 100
)

type Article struct {
	Title   string `json: Title`
	Desc    string `json:desc`
	Content string `json:content`
}

type Articles []Article

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

func getCharacters(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Description", Content: "Hello World"},
	}
	fmt.Println("Endpoint Hit: All Articles Endpoint")
	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		return
	}
}

func getCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	err := json.NewEncoder(w).Encode(&Book{})
	if err != nil {
		return 
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/characters", getCharacters).Methods("GET")
	router.HandleFunc("/characters/{id}", getCharacter).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type config struct {
	privateKey string
	publicKey string
}

// TODO: Change to read from config file
func readConfig() *config {
	return &config{
		privateKey: os.Getenv("MARVEL_API_PRIVATE_KEY"),
		publicKey: os.Getenv("MARVEL_API_PUBLIC_KEY"),
	}
}

func main() {
	// Init Router
	//router := mux.NewRouter()
	conf := readConfig()
	fmt.Fprintln(os.Stderr, conf)
	handleRequests()
}
