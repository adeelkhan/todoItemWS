package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// models
type TodoItem struct {
	Id string `json:"Id"`
	ItemName string `json:"item_name"` 
	CreateTimeStamp time.Time `json:"create_timestamp"`
	UpdateTimeStamp time.Time `json:"update_timestamp"`
}

// request/response structures
type ItemCreateRequest struct {
	ItemName string  `json:"item_name"`
}

type ItemDeleteRequest struct {
	Id string `json:"item_id"`
}

type ItemUpdateRequest struct {
	Id string `json:"item_id"`
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

type LoginResponse struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
	User string `json:"user"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserProfile struct {
	UserName string
	Password string
	todoItem map[string]string
}

var jwtKey = []byte("")

// model
var todoMap = map[string]TodoItem{}

var users = map[string]*UserProfile {
	"user1@abc.com": {
		UserName: "user1@abc.com",
		Password: "password1",
	},
	"user2@abc.com": {
		UserName: "user2@abc.com",
		Password: "password2",
	},
}

func addItem(username, itemId string, todo *TodoItem) {
	todoMap[itemId] = *todo
	if users[username].todoItem == nil {
		users[username].todoItem = make(map[string]string, 1)
	}
	users[username].todoItem[itemId] = itemId
}

func removeItem(username, itemId string) {
	delete(todoMap, itemId)
	delete(users[username].todoItem, itemId)
}

// handlers
func createItem(w http.ResponseWriter, req *http.Request){
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}
	// check jwt validity
	username := getUser(w, req)
	decode := json.NewDecoder(req.Body)
	
	request := ItemCreateRequest{}
	err := decode.Decode(&request)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	itemId := uuid.NewString()
	createTime := time.Now() 
	var todo = TodoItem{
		Id: itemId,
		ItemName: request.ItemName,
		CreateTimeStamp: createTime,
		UpdateTimeStamp: createTime,
	}
	fmt.Println(todo)
	addItem(username, itemId, &todo)

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
	// check jwt validity
	username := getUser(w, req)

	decode := json.NewDecoder(req.Body)
	request := ItemDeleteRequest{}
	err := decode.Decode(&request)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	removeItem(username, request.Id)
	delete(todoMap, request.Id)

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

	var Id string
	if Id != "" {
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
	// check jwt validity
	username := getUser(w, req)
	userProfile := users[username]


	items := make([]TodoItem,0)

	for _, itemsId := range userProfile.todoItem {
		v, ok := todoMap[itemsId]
		if ok {
			todo := TodoItem{
				Id: v.Id,
				ItemName: v.ItemName,
				CreateTimeStamp: v.CreateTimeStamp,
				UpdateTimeStamp: v.UpdateTimeStamp,
			}
			items = append(items, todo)
		}
	}

	var response = ListItemResponse{
		Msg: "Success",
		Status: http.StatusOK,
		Items: items,
	}

	res, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", string(res))
}

func Signin(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	var creds Credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	userProfile, ok := users[creds.Username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if creds.Password != userProfile.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return 
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := Claims{
		Username: creds.Username, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	http.SetCookie(w, &http.Cookie{
		Name: "token", 
		Value: tokenString, 
		Expires: expirationTime,
	})

	var response = LoginResponse{
		Msg: "Success",
		User: creds.Username,
		Status: http.StatusOK,
	}

	res, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", string(res))

}

func Refresh(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	c, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)		
			return 
		}
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	tknStr := c.Value 
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token)(any, error) {
		return jwtKey, nil 
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}
		w.WriteHeader(http.StatusBadRequest)
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return 
	}

	if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	// now, create a new token for the current use, wit a renewed expiration 
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	// set the new token as the users `token` cookie 
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString, 
		Expires: expirationTime,
	})
	
}
func Logout(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if (*req).Method == "OPTIONS" {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Expires: time.Now(),
	})
}

// check jwt validity 
func getUser(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}
	tknStr := c.Value 
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil 
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return ""
	}
	return claims.Username
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
func main() {
	
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", Logout)

	http.HandleFunc("/create", createItem)
	http.HandleFunc("/delete", deleteItem)
	http.HandleFunc("/update", updateItem)
	http.HandleFunc("/list", listItem)
	
	fmt.Println("Server started on port 8090")
	http.ListenAndServe(":8090", nil)
}