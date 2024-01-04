package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cred = Credentials{
	Username: "user1@abc.com",
	Password: "password1",
}

var testItemName = "test-item"

// setting up server engine
var router = setupRouter()

func GetAuthToken() []string {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(cred)

	request, _ := http.NewRequest("POST", "/signin", buf)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	token := response.Header()["Set-Cookie"]

	return token
}

func TestAuth(t *testing.T) {
	// dont want test to emit any log
	logger = slog.New(slog.NewJSONHandler(ioutil.Discard, nil))

	// setting up payload
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(cred)

	request, _ := http.NewRequest("POST", "/signin", buf)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	lr := LoginResponse{}
	err := json.NewDecoder(response.Body).Decode(&lr)
	if err != nil {
		fmt.Println(err.Error())
	}
	want := LoginResponse{
		Msg:    "Success",
		Status: 200,
		User:   cred.Username,
	}
	if lr != want {
		t.Errorf("got %v, want %v", lr, want)
	}
}

// func TestWS
func TestCreateItem(t *testing.T) {
	item := ItemCreateRequest{
		ItemName: testItemName,
	}

	// setting up request
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(item)
	request, _ := http.NewRequest("POST", "/create", buf)
	response := httptest.NewRecorder()

	// setting up the cookie
	token := GetAuthToken()
	request.Header["Cookie"] = token

	router.ServeHTTP(response, request)

	ir := ItemResponse{}
	err := json.NewDecoder(response.Body).Decode(&ir)
	if err != nil {
		fmt.Println(err.Error())
	}
	want := ItemResponse{
		Msg:    "Success",
		Status: 201,
	}
	if ir != want {
		t.Errorf("got %v, want %v", ir, want)
	}
}

func TestGetItemList(t *testing.T) {
	request, _ := http.NewRequest("GET", "/list", nil)
	response := httptest.NewRecorder()

	token := GetAuthToken()

	// setting up the cookie
	request.Header["Cookie"] = token
	router.ServeHTTP(response, request)

	items := ListItemResponse{}

	err := json.NewDecoder(response.Body).Decode(&items)
	if err != nil {
		fmt.Println(err.Error())
	}
	want := 1
	if len(items.Items) != 1 {
		t.Errorf("got %v, want %v", len(items.Items), want)
	}
}

func TestUpdateItem(t *testing.T) {
	userProfile := users[cred.Username]
	todoItemsUUIDs := userProfile.todoItem

	todoIds := make([]string, 0)
	for id, _ := range todoItemsUUIDs {
		todoIds = append(todoIds, id)
	}

	itemUpdated := ItemUpdateRequest{Id: todoIds[0], ItemName: "UpdatedItem"}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(itemUpdated)

	request, _ := http.NewRequest("POST", "/update", buf)
	response := httptest.NewRecorder()

	token := GetAuthToken()
	request.Header["Cookie"] = token
	router.ServeHTTP(response, request)

	ir := ItemResponse{}
	err := json.NewDecoder(response.Body).Decode(&ir)
	if err != nil {
		fmt.Println(err.Error())
	}
	want := ItemResponse{
		Msg:    "Success",
		Status: 200,
	}
	if ir != want {
		t.Errorf("got %v, want %v", ir, want)
	}
}

func TestDeleteItem(t *testing.T) {

	userProfile := users[cred.Username]
	todoItemsUUIDs := userProfile.todoItem

	todoIds := make([]string, 0)
	for id, _ := range todoItemsUUIDs {
		todoIds = append(todoIds, id)
	}

	itemDelete := ItemDeleteRequest{Id: todoIds[0]}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(itemDelete)

	request, _ := http.NewRequest("POST", "/delete", buf)
	response := httptest.NewRecorder()

	// setup cookie
	token := GetAuthToken()
	request.Header["Cookie"] = token
	router.ServeHTTP(response, request)

	ir := ItemResponse{}
	err := json.NewDecoder(response.Body).Decode(&ir)
	if err != nil {
		fmt.Println(err.Error())
	}
	want := ItemResponse{
		Msg:    "Success",
		Status: 200,
	}
	if ir != want {
		t.Errorf("got %v, want %v", ir, want)
	}
}
