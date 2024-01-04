package main

import (
	"net/http"
	"os"
	"sort"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	// "github.com/rs/cors"
	cors "github.com/rs/cors/wrapper/gin"
)

var handler = slog.NewJSONHandler(os.Stdout, nil)
var logger = slog.New(handler)

// models
type TodoItem struct {
	Id              string    `json:"Id"`
	ItemName        string    `json:"item_name"`
	CreateTimeStamp time.Time `json:"create_timestamp"`
	UpdateTimeStamp time.Time `json:"update_timestamp"`
}

// request/response structures
type ItemCreateRequest struct {
	ItemName string `json:"item_name"`
}

type ItemDeleteRequest struct {
	Id string `json:"item_id"`
}

type ItemUpdateRequest struct {
	Id       string `json:"item_id"`
	ItemName string `json:"item_name"`
}

type ItemResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type ListItemResponse struct {
	Msg    string     `json:"msg"`
	Status int        `json:"status"`
	Items  []TodoItem `json:"items"`
}

type LoginResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	User   string `json:"user"`
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

var users = map[string]*UserProfile{
	"user1@abc.com": {
		UserName: "user1@abc.com",
		Password: "password1",
	},
	"user2@abc.com": {
		UserName: "user2@abc.com",
		Password: "password2",
	},
}

func AddItem(username, itemId string, todo *TodoItem) {
	todoMap[itemId] = *todo
	if users[username].todoItem == nil {
		users[username].todoItem = make(map[string]string, 1)
	}
	users[username].todoItem[itemId] = itemId
}

func RemoveItem(username, itemId string) {
	delete(todoMap, itemId)
	delete(users[username].todoItem, itemId)
}

// handlers
func CreateItem(c *gin.Context) {

	// check jwt validity
	username := getUser(c)

	request := ItemCreateRequest{}
	if err := c.BindJSON(&request); err != nil {
		logger.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, check payload"})
		return
	}

	itemId := uuid.NewString()
	createTime := time.Now()
	var todo = TodoItem{
		Id:              itemId,
		ItemName:        request.ItemName,
		CreateTimeStamp: createTime,
		UpdateTimeStamp: createTime,
	}
	AddItem(username, itemId, &todo)

	logger.Info("Item created successfully",
		slog.String("Item Id#", todo.Id),
		slog.String("item Name", todo.ItemName),
		slog.Time("Created At", todo.CreateTimeStamp),
		slog.Time("Updated At", todo.UpdateTimeStamp),
	)

	res := ItemResponse{
		Msg:    "Success",
		Status: http.StatusCreated,
	}
	c.IndentedJSON(http.StatusCreated, res)
}

func DeleteItem(ctx *gin.Context) {
	request := ItemDeleteRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, check payload"})
		return
	}

	// check jwt validity
	username := getUser(ctx)

	RemoveItem(username, request.Id)
	delete(todoMap, request.Id)
	logger.Info("Item deleted successfully",
		slog.String("Item Id#", request.Id),
	)

	res := ItemResponse{
		Msg:    "Success",
		Status: http.StatusOK,
	}
	ctx.IndentedJSON(http.StatusOK, res)
}
func UpdateItem(ctx *gin.Context) {
	request := ItemUpdateRequest{}
	if err := ctx.BindJSON(&request); err != nil {
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, check payload"})
	}

	var Id string = request.Id
	if Id != "" {
		todoMap[Id] = TodoItem{
			Id:              Id,
			ItemName:        request.ItemName,
			CreateTimeStamp: todoMap[Id].CreateTimeStamp,
			UpdateTimeStamp: time.Now(),
		}
	}
	logger.Info("Item updated successfully",
		slog.String("Item Id#", request.Id),
	)

	res := ItemResponse{
		Msg:    "Success",
		Status: http.StatusOK,
	}
	ctx.IndentedJSON(http.StatusOK, res)
}

func ListItem(ctx *gin.Context) {
	// check jwt validity
	username := getUser(ctx)
	userProfile := users[username]

	items := make([]TodoItem, 0)

	for _, itemsId := range userProfile.todoItem {
		v, ok := todoMap[itemsId]
		if ok {
			todo := TodoItem{
				Id:              v.Id,
				ItemName:        v.ItemName,
				CreateTimeStamp: v.CreateTimeStamp,
				UpdateTimeStamp: v.UpdateTimeStamp,
			}
			items = append(items, todo)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreateTimeStamp.Unix() < items[j].CreateTimeStamp.Unix()
	})

	var response = ListItemResponse{
		Msg:    "Success",
		Status: http.StatusOK,
		Items:  items,
	}
	ctx.IndentedJSON(http.StatusOK, response)
}

func Signin(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		logger.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, check your input"})
	}
	userProfile, ok := users[creds.Username]
	if !ok {
		logger.Error("Invalid Credentials, requested user not found")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Username/password is incorrect"})
	}

	if creds.Password != userProfile.Password {
		logger.Error("Invalid credentials passed")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Username/password is incorrect"})
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logger.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something unexpected happened, try again."})
	}
	c.SetCookie("token", tokenString, 60*30, "/", "", false, false)

	response := LoginResponse{
		Msg:    "Success",
		User:   creds.Username,
		Status: http.StatusOK,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func Refresh(ctx *gin.Context) {
	tknStr, err := ctx.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			logger.Error(err.Error())
			ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
		}
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, invalid jwt"})
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error(err.Error())
			ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
		}
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, invalid jwt"})
	}
	if !tkn.Valid {
		logger.Error("Invalid Jwt token")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
	}

	if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
		logger.Error("Token refresh request is too close to expiry time")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, invalid jwt"})
	}
	// now, create a new token for the current use, wit a renewed expiration
	expirationTime := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Something unexpected happened, try again."})
	}
	ctx.SetCookie("token", tokenString, expirationTime.Minute(), "/", "", false, false)

}
func Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", 0, "/", "", false, false)
}

//  check jwt validity
func getUser(ctx *gin.Context) string {
	tknStr, err := ctx.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			logger.Error(err.Error())
			ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
			return ""
		}
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, invalid jwt"})
		return ""
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error(err.Error())
			ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
			return ""
		}
		logger.Error(err.Error())
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad request, invalid jwt"})
		return ""
	}
	if !tkn.Valid {
		logger.Error("Invalid Jwt token")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Jwt"})
		return ""
	}
	return claims.Username
}
func setupRouter() *gin.Engine {
	router := gin.Default()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"PUT", "PATCH", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true,
	})
	router.Use(c)
	router.POST("/signin", Signin)
	router.POST("/refresh", Refresh)
	router.POST("/create", CreateItem)
	router.POST("/update", UpdateItem)
	router.POST("/delete", DeleteItem)
	router.GET("/list", ListItem)

	return router
}
func main() {
	router := setupRouter()
	logger.Info("Server starting on port 8091")
	router.Run("localhost:8091")
}
