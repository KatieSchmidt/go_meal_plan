package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
	// "reflect"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meal struct {
	Username string  `json:"username" bson:"username"`
	Mealname  string `json:"mealname" bson:"mealname"`
	TotalCalories int64 `json:"totalcalories" bson:"totalcalories"`
	Ingredients []Ingredient `json:"ingredients" bson:"ingredients"`
	DateAdded time.Time
}

type Ingredient struct {
	ID primitive.ObjectID `json:"ing_id" bson:"ing_id"`
	Ingredient string `json:"ingredient" bson:"ingredient"`
	Calories int64 `json:"calories" bson:"calories"`
	MeasureUnitQuantity int64 `json:"measureunitquantity" bson:"measureunitquantity"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
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
		meal.DateAdded = time.Now()
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

func GetMealById(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	collection := client.Database("go_meals").Collection("meals")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//make meal struc and get/make objectID, find by that id
	var resulting_meal Meal
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id }}
	error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
	//if the meal wasnt found, create it else send an error message
	if error_msg != nil {
		response_message := ErrorMessage{"meal not found"}
		json.NewEncoder(response).Encode(response_message)
	} else {
		json.NewEncoder(response).Encode(resulting_meal)
	}
}

func AddIngredientToMeal(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	response.Header().Set("content-type", "application/x-www-form-urlencoded")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id }}

	var ingredient Ingredient
	ingredient.ID = primitive.NewObjectID()
	ingredient.Ingredient = request.FormValue("ingredient")
	if cals, err := strconv.ParseInt(request.FormValue("calories"), 10, 32); err == nil {
		ingredient.Calories = cals
	}
	ingredient.MeasureUnit = request.FormValue("measureunit")
	if units, err := strconv.ParseInt(request.FormValue("measureunitquantity"), 10, 32); err == nil {
		ingredient.MeasureUnitQuantity = units
	}

	collection := client.Database("go_meals").Collection("meals")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var resulting_meal Meal
	var updated_meal Meal
	error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
	if error_msg != nil {
		response_message := ErrorMessage{"meal not found"}
		json.NewEncoder(response).Encode(response_message)
	} else {
		resulting_meal.Ingredients = append(resulting_meal.Ingredients, ingredient)
		resulting_meal.TotalCalories = resulting_meal.TotalCalories + ingredient.Calories
		error_msg_2 := collection.FindOneAndReplace(ctx, filter, resulting_meal).Decode(&updated_meal)

		if error_msg_2 != nil  {
			response_message := ErrorMessage{"Unable to add ingredient"}
			json.NewEncoder(response).Encode(response_message)
		} else {
			json.NewEncoder(response).Encode(resulting_meal)
		}
	}
}


func DeleteMealById(response http.ResponseWriter, request *http.Request){
		params := mux.Vars(request)
		collection := client.Database("go_meals").Collection("meals")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//make objectID, find by that id
		id, _ := primitive.ObjectIDFromHex(params["id"])
		filter := bson.D{{"_id", id }}
		result, error_msg := collection.DeleteOne(ctx, filter)
		//if the meal wasnt found, create it else send an error message
		if error_msg != nil {
			response_message := ErrorMessage{"meal not found"}
			json.NewEncoder(response).Encode(response_message)
		} else {
			json.NewEncoder(response).Encode(result)
		}

}

func DeleteIngredientFromMeal(response http.ResponseWriter, request *http.Request){
	fmt.Println("This will delete an ingredient from meal by meal id")
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
	router.HandleFunc("/meals/{id}", GetMealById).Methods("GET")
	router.HandleFunc("/meals/{id}/ingredients", AddIngredientToMeal).Methods("PUT")
	router.HandleFunc("/meals/{id}", DeleteMealById).Methods("DELETE")
	router.HandleFunc("/meals/{id}/ingredients/{ingredient}", DeleteIngredientFromMeal).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":5000", router))
}
