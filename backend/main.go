package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// models
type TodoItem struct {
	Id int `json:"Id"`
	ItemName string `json:"item_name"` 
	CreateTimeStamp time.Time `json:"create_timestamp"`
	UpdateTimeStamp time.Time `json:"update_timestamp"`
}

// request/response structures
type ItemCreateRequest struct {
	ItemName string  `json:"item_name"`
}

type ItemDeleteRequest struct {
	Id int `json:"item_id"`
}

type ItemUpdateRequest struct {
	Id int `json:"item_id"`
	ItemName string `json:"item_name"`
}

type ItemResponse struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
}

type ListItemResponse struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
	Items []TodoItem `json:"items"`
}


// model
var todoMap = map[int]TodoItem{}

// handlers
func createItem(w http.ResponseWriter, req *http.Request){
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}
	
	decode := json.NewDecoder(req.Body)
	
	request := ItemCreateRequest{}
	err := decode.Decode(&request)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	createTime := time.Now() 
	var todo = TodoItem{
		Id: len(todoMap) + 1,
		ItemName: request.ItemName,
		CreateTimeStamp: createTime,
		UpdateTimeStamp: createTime,
	}

	fmt.Println(todo)

	todoMap[todo.Id] = todo
	res ,_ := json.Marshal(ItemResponse{
		Msg: "Success",
		Status: http.StatusCreated,
	})
	fmt.Fprintf(w, "%s", string(res))
}

func deleteItem(w http.ResponseWriter, req *http.Request){
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	decode := json.NewDecoder(req.Body)
	request := ItemDeleteRequest{}
	err := decode.Decode(&request)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	for _,v := range todoMap {
		if request.Id == v.Id {
			delete(todoMap, v.Id)
			break
		}
	}

	res ,_ := json.Marshal(ItemResponse{
		Msg: "Success",
		Status: http.StatusOK,
	})
	fmt.Fprintf(w, "%s", string(res))
}
func updateItem(w http.ResponseWriter, req *http.Request){
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	decode := json.NewDecoder(req.Body)
	request := ItemUpdateRequest{}
	err := decode.Decode(&request)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	var Id int
	for _,v := range todoMap {
		if request.Id == v.Id {
			Id = v.Id
			break
		}
	}
	if Id != 0 {
		todoMap[Id] = TodoItem{
			Id: Id,
			ItemName: request.ItemName,
			CreateTimeStamp: todoMap[Id].CreateTimeStamp,
			UpdateTimeStamp: time.Now(),
		}
	}
	res ,_ := json.Marshal(ItemResponse{
		Msg: "Success",
		Status: http.StatusOK,
	})
	fmt.Fprintf(w, "%s", string(res))
}
func listItem(w http.ResponseWriter, req *http.Request){
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	items := make([]TodoItem,0)
	
	for _, v := range todoMap {
		todo := TodoItem{
			Id: v.Id,
			ItemName: v.ItemName,
			CreateTimeStamp: v.CreateTimeStamp,
			UpdateTimeStamp: v.UpdateTimeStamp,
		}
		items = append(items, todo)
	}

	var response = ListItemResponse{
		Msg: "Success",
		Status: http.StatusOK,
		Items: items,
	}

	res, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", string(res))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
func main() {
	http.HandleFunc("/create", createItem)
	http.HandleFunc("/delete", deleteItem)
	http.HandleFunc("/update", updateItem)
	http.HandleFunc("/list", listItem)
	
	http.ListenAndServe(":8090", nil)
}