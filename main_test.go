package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)
}

func TestCharactersResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/characters", nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got <%v> want <%v>", status, http.StatusOK)
	}
}

func TestCharactersInvalidId(t *testing.T) {
	req, err := http.NewRequest("GET", "/characters/123", nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if code := rr.Code; code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got <%v> want <%v>", code, http.StatusOK)
	}

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	expectedStatus := "We couldn't find that character"
	if m["status"] != expectedStatus {
		t.Errorf("handler returned wrong status: got <%v> want <%v>", m["status"], expectedStatus)
	}
}

func TestCharactersValidId(t *testing.T) {
	req, err := http.NewRequest("GET", "/characters/1010698", nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got <%v> want <%v>", status, http.StatusOK)
	}

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	expectedId := float64(1010698)
	if m["id"] != expectedId {
		t.Errorf("handler returned wrong id: got <%v> want <%v>", m["id"], expectedId)
	}

	expectedName := "Young Avengers"
	if m["name"] != expectedName {
		t.Errorf("handler returned wrong name: got <%v> want <%v>", m["name"], expectedName)
	}
}
