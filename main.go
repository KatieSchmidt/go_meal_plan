package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	// "reflect"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meal struct {
	Username string  `json:"username" bson:"username"`
	Mealname  string `json:"mealname" bson:"mealname"`
}

type ErrorMessage struct {
	Error string
}

func CreateMeal(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	response.Header().Set("content-type", "application/x-www-form-urlencoded")

	if len(request.FormValue("username")) == 0 || len(request.FormValue("mealname")) == 0{
		meal_error := ErrorMessage{"One of your form fields was empty"}
		json.NewEncoder(response).Encode(meal_error)
	} else {
		// look for a meal that has same username and meal name
			// if it does return an error if not, make the meal
		collection := client.Database("go_meals").Collection("meals")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var meal Meal
		meal.Username = request.FormValue("username")
		meal.Mealname = request.FormValue("mealname")
		filter := bson.D{{"username", meal.Username}, {"mealname", meal.Mealname}}
		var resulting_meal Meal
		error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
		//if the meal wasnt found, create it else send an error message
		if error_msg != nil {
			//create Meal using the form data, save to a collection,
			_, err := collection.InsertOne(ctx, meal)
			if err != nil {
				response_message := ErrorMessage{"ERROR: there was an error creating your meal"}
				json.NewEncoder(response).Encode(response_message)
			} else {
				//if there isnt an error, meal was inserted, so return the meal
				json.NewEncoder(response).Encode(meal)
			}
		} else {
			response_message := ErrorMessage{"meal already exists"}
			json.NewEncoder(response).Encode(response_message)
		}
	}
}

func GetMeals(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := client.Database("go_meals").Collection("meals")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
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

	if len(meals) > 0 {
		json.NewEncoder(response).Encode(meals)

	} else {
		//if there are no meals create a message Struct to send back
		response_message := ErrorMessage{"Error: No meals have been created"}
		json.NewEncoder(response).Encode(response_message)
	}
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
