package meals

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "github.com/KatieSchmidt/meal_plan/models"
)

func GetMeals(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func (response http.ResponseWriter, request *http.Request) {
  	response.Header().Set("content-type", "application/json")
  	collection := mongoClient.Database("go_meals").Collection("meals")
  	cursor, err := collection.Find(ctx, bson.M{})

  	if err != nil {
  		log.Fatal(err)
  	}

  	//create a list of meals of struc models.Meal
  	var meals []models.Meal
  	for cursor.Next(ctx) {
  		var meal models.Meal
  		cursor.Decode(&meal)
  		meals = append(meals, meal)
  	}

  	if len(meals) > 0 {
  		json.NewEncoder(response).Encode(meals)

  	} else {
  		//if there are no meals create a message Struct to send back
  		response_message := models.ErrorMessage{"Error: No meals have been created"}
  		json.NewEncoder(response).Encode(response_message)
  	}
  }
}

func CreateMeal(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request) {
  	request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")

  	if len(request.FormValue("user")) == 0 || len(request.FormValue("mealname")) == 0{
  		meal_error := models.ErrorMessage{"One of your form fields was empty"}
  		json.NewEncoder(response).Encode(meal_error)
  	} else {
  		// look for a meal that has same username and meal name
  			// if it does return an error if not, make the meal
  		collection := mongoClient.Database("go_meals").Collection("meals")
  		var meal models.Meal

	  	meal.ID = primitive.NewObjectID()
  		meal.User, _ = primitive.ObjectIDFromHex(request.FormValue("user"))
  		meal.Mealname = request.FormValue("mealname")
  		meal.DateAdded = time.Now()
  		filter := bson.D{{"user", meal.User}, {"mealname", meal.Mealname}}
  		var resulting_meal models.Meal
  		error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
  		//if the meal wasnt found, create it else send an error message
  		if error_msg != nil {
  			//create models.Meal using the form data, save to a collection,
  			_, err := collection.InsertOne(ctx, meal)
  			if err != nil {
  				response_message := models.ErrorMessage{"ERROR: there was an error creating your meal"}
  				json.NewEncoder(response).Encode(response_message)
  			} else {
  				//if there isnt an error, meal was inserted, so return the meal
  				json.NewEncoder(response).Encode(meal)
  			}
  		} else {
  			response_message := models.ErrorMessage{"meal already exists"}
  			json.NewEncoder(response).Encode(response_message)
  		}
  	}
  }
}

func GetMealById(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request) {
  	params := mux.Vars(request)
  	collection := mongoClient.Database("go_meals").Collection("meals")
  	//make meal struc and get/make objectID, find by that id
  	var resulting_meal models.Meal
  	meal_id, _ := primitive.ObjectIDFromHex(params["meal_id"])
  	filter := bson.D{{"_id", meal_id }}
  	error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
  	//if the meal wasnt found, create it else send an error message
  	if error_msg != nil {
  		response_message := models.ErrorMessage{"meal not found"}
  		json.NewEncoder(response).Encode(response_message)
  	} else {
  		json.NewEncoder(response).Encode(resulting_meal)
  	}
  }
}

func AddIngredientToMeal(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request) {
  	request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")

  	params := mux.Vars(request)
  	meal_id, _ := primitive.ObjectIDFromHex(params["meal_id"])
  	filter := bson.D{{"_id", meal_id }}

  	var ingredient models.Ingredient
  	ingredient.ID = primitive.NewObjectID()
  	ingredient.Ingredient = request.FormValue("ingredient")
  	if cals, err := strconv.ParseFloat(request.FormValue("calories"), 64); err == nil {
  		ingredient.Calories = cals
  	}
  	ingredient.MeasureUnit = request.FormValue("measureunit")
  	if units, err := strconv.ParseFloat(request.FormValue("measureunitquantity"), 64); err == nil {
  		ingredient.MeasureUnitQuantity = units
  	}

  	collection := mongoClient.Database("go_meals").Collection("meals")

  	var resulting_meal models.Meal
  	var updated_meal models.Meal
  	error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
  	if error_msg != nil {
  		response_message := models.ErrorMessage{"meal not found"}
  		json.NewEncoder(response).Encode(response_message)
  	} else {
  		resulting_meal.Ingredients = append(resulting_meal.Ingredients, ingredient)
  		resulting_meal.TotalCalories = resulting_meal.TotalCalories + ingredient.Calories
  		error_msg_2 := collection.FindOneAndReplace(ctx, filter, resulting_meal).Decode(&updated_meal)

  		if error_msg_2 != nil  {
  			response_message := models.ErrorMessage{"Unable to add ingredient"}
  			json.NewEncoder(response).Encode(response_message)
  		} else {
  			json.NewEncoder(response).Encode(resulting_meal)
  		}
  	}
  }
}
func DeleteMealById(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request){
  		params := mux.Vars(request)
  		collection := mongoClient.Database("go_meals").Collection("meals")
  		//make objectID, find by that id
  		id, _ := primitive.ObjectIDFromHex(params["meal_id"])
  		filter := bson.D{{"_id", id }}
  		result, _ := collection.DeleteOne(ctx, filter)
  		//if the meal wasnt found, create it else send an error message
			// DeleteOne always returns a result. error is always nil. so check to see if deleted count is equal to zero instead
	    if result.DeletedCount == 0 {
	      response_message := models.ErrorMessage{"meal not found"}
	      json.NewEncoder(response).Encode(response_message)
	    } else {
	      response_message := models.ErrorMessage{"meal deleted"}
	      json.NewEncoder(response).Encode(response_message)
	    }
  }
}
func DeleteIngredientFromMeal(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
  return func(response http.ResponseWriter, request *http.Request){
  	params := mux.Vars(request)
  	collection := mongoClient.Database("go_meals").Collection("meals")
  	//make meal struc and get/make objectID, find by that id
  	var resulting_meal models.Meal
  	meal_id, _ := primitive.ObjectIDFromHex(params["meal_id"])
  	ing_id, _ := primitive.ObjectIDFromHex(params["ingredient_id"])
  	filter := bson.D{{"_id", meal_id }}
  	error_msg := collection.FindOne(ctx, filter).Decode(&resulting_meal)
  	//find meal
  	//ingredients is a slice. find index of ingredient to know where to remove

  	if error_msg != nil {
  		response_message := models.ErrorMessage{"meal not found"}
  		json.NewEncoder(response).Encode(response_message)
  	} else {
  		ingredients := resulting_meal.Ingredients
  		index := 0
  		running := true
  		var calories_to_remove float64 = 0
  		for i := 0; i < len(ingredients) && running == true; i++ {
  			if ingredients[i].ID == ing_id {
  				index = i
  				calories_to_remove = ingredients[i].Calories
  				running = false
  			}
  		}
  		new_ing_slice := append(ingredients[:index], ingredients[index+1:]...)
  		fmt.Println(new_ing_slice)

  		resulting_meal.Ingredients = new_ing_slice
  		resulting_meal.TotalCalories = resulting_meal.TotalCalories - calories_to_remove

  		var updated_meal models.Meal
  		error_msg_2 := collection.FindOneAndReplace(ctx, filter, resulting_meal).Decode(&updated_meal)

  		if error_msg_2 != nil  {
  			response_message := models.ErrorMessage{"Unable to remove ingredient"}
  			json.NewEncoder(response).Encode(response_message)
  		} else {
  			json.NewEncoder(response).Encode(resulting_meal)
  		}
  	}
  }
}
