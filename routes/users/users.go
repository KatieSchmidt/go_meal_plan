package users

import (
  "fmt"
	"encoding/json"
	"context"
  "time"
  // "strings"
  // "log"
  // "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "net/http"
	"github.com/KatieSchmidt/meal_plan/models"
  "golang.org/x/crypto/bcrypt"
  "github.com/dgrijalva/jwt-go"
)

// POST Registers user encrypts password
func RegisterUser(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request) {
    fmt.Println("working")
    request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")
    collection := mongoClient.Database("go_meals").Collection("users")
    filter := bson.D{{"email", request.FormValue("email")}}

    var currentUser models.User
    err := collection.FindOne(ctx, filter).Decode(&currentUser)

    if err != nil {
      bstring := []byte(request.FormValue("password"))
      bcryptPassword, _ := bcrypt.GenerateFromPassword(bstring, 10)

      var newUser models.User
      newUser.ID = primitive.NewObjectID()
      newUser.Name = request.FormValue("name")
      newUser.Email = request.FormValue("email")
      newUser.Password = bcryptPassword

      _, err := collection.InsertOne(ctx, newUser)

      if err != nil {
        var errorMessage models.Errors
        errorMessage.User = "there was an error creating the user"
        json.NewEncoder(response).Encode(errorMessage)
      } else {
        json.NewEncoder(response).Encode(newUser)
      }
    } else {
      var errorMessage models.Errors
      errorMessage.User = "This email already exists for a user."
      json.NewEncoder(response).Encode(errorMessage)
    }
  }
}

//POST logs in user/returns a JWT
func LoginUser(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request) {
    request.ParseForm()
    response.Header().Set("content-type", "application/x-www-form-urlencoded")
    collection := mongoClient.Database("go_meals").Collection("users")
    email := request.FormValue("email")
    password := []byte(request.FormValue("password"))

    filter := bson.D{{"email", email}}
    var user models.User
    err := collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
      var errMessage models.Errors
      errMessage.User = "This user not in database"
      json.NewEncoder(response).Encode(errMessage)
    } else {
      bcryptErr := bcrypt.CompareHashAndPassword(user.Password, password)
      if bcryptErr != nil {
        var errMessage models.Errors
        errMessage.Password = "Password email error"
        json.NewEncoder(response).Encode(errMessage)
        response.WriteHeader(http.StatusUnauthorized)
      } else {
        expirationTime := time.Now().Add(1440 * time.Minute)
        claims := models.Claims{
        		Name: user.Name,
            ID: user.ID,
        		StandardClaims: jwt.StandardClaims{
        			ExpiresAt: expirationTime.Unix(), //so its ms
        		},
        	}
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        var jwtKey = []byte("my_secret_key")
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
          var errMessage models.Errors
          errMessage.Password = "There was an error signing your token"
          json.NewEncoder(response).Encode(errMessage)
          response.WriteHeader(http.StatusUnauthorized)
        } else {
          var JOT models.JOT
          JOT.Success = true
          JOT.Token = tokenString
          json.NewEncoder(response).Encode(JOT)
        }
      }
    }
  }
}
