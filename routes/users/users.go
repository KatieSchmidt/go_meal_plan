package users

import (
  "fmt"
	"encoding/json"
	"context"
  // "strings"
  // "log"
  // "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "net/http"
	"github.com/KatieSchmidt/meal_plan/models"
  "golang.org/x/crypto/bcrypt"
  // "github.com/dgrijalva/jwt-go"
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

  }
}
