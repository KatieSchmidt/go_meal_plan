package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"reflect"
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

func CreateMeal(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	response.Header().Set("content-type", "application/x-www-form-urlencoded")

	//create Meal using the form data, save to a collection,
	var meal Meal
	meal.Username = request.FormValue("username")
	meal.Mealname = request.FormValue("mealname")
	collection := client.Database("go_meals").Collection("meals")
	//cancel will cancel ctx as soon as timeout expires
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("meal: ", meal, reflect.TypeOf(meal))
	result, err := collection.InsertOne(ctx, meal)

	if err != nil {
		type ErrorMessage struct {
			Error string
		}
		response_message := ErrorMessage{"ERROR: there was an error creating your meal"}
		json.NewEncoder(response).Encode(response_message)
	} else {
		fmt.Println(result, reflect.TypeOf(result))
		// filter := bson.D{"_id", result}
		single_result := collection.FindOne(ctx, result)
		fmt.Println("single_result: ", single_result, reflect.TypeOf(single_result))

		var new_meal Meal
		single_result.Decode(&new_meal)
		fmt.Println("new_meal: ", new_meal, reflect.TypeOf(new_meal))
		json.NewEncoder(response).Encode(new_meal)
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
		type ErrorMessage struct {
			Error string
		}
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
