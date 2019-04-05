package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meal struct {
	Username string  `json:"username,omitempty" bson:"username,omitempty"`
	Mealname  string `json:"mealname,omitempty" bson:"mealname,omitempty"`
}

func CreateMeal(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	response.Header().Set("content-type", "application/x-www-form-urlencoded")

	//ctreate Meal using the form data, save to a collection,
	var meal Meal
	meal.Username = request.FormValue("username")
	meal.Mealname = request.FormValue("mealname")
	collection := client.Database("meal_plan").Collection("meals")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meal)
	json.NewEncoder(response).Encode(result)
}

func GetMeals(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("go_meals").Collection("meals")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("ERROR: no meals were found!"))
		log.Fatal(err)

	}
	defer cursor.Close(ctx)


	//create a list of meals of struc Meal
	var meals []Meal
	for cursor.Next(ctx) {
		var meal Meal
		cursor.Decode(&meal)
		meals = append(meals, meal)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(response).Encode(meals)
}


func main() {
	fmt.Println("Server started on port 5000")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//try to do the rest of everything and then cancel
	defer cancel()

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/meal", CreateMeal).Methods("POST")
	router.HandleFunc("/meals", GetMeals).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}
