package mealplans

import (
	"fmt"
	"encoding/json"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
  "net/http"
	"github.com/KatieSchmidt/meal_plan/models"
)


func  CreateMealplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")

  	if len(request.FormValue("userid")) == 0 || len(request.FormValue("planname")) == 0 {
  		meal_error := models.ErrorMessage{"One of your form fields was empty"}
  		json.NewEncoder(response).Encode(meal_error)
  	} else {
  		// look for a mealplan that has same name and user
  			// if it does return an error if not, make the meal
  		collection := mongoClient.Database("go_meals").Collection("mealplans")
  		var mealplan models.Mealplan
  		mealplan.UserId = request.FormValue("userid")
  		mealplan.Planname = request.FormValue("planname")
  		filter := bson.D{{"userid", mealplan.UserId}, {"planname", mealplan.Planname}}
  		error_msg := collection.FindOne(ctx, filter)

  		if error_msg != nil {
  			_, err := collection.InsertOne(ctx, mealplan)
  			if err != nil {
  				response_message := models.ErrorMessage{"ERROR: there was an error creating your mealplan"}
  				json.NewEncoder(response).Encode(response_message)
  			} else {
  				//if there isnt an error, meal was inserted, so return the meal
  				json.NewEncoder(response).Encode(mealplan)
  			}
  		} else {
  			error := models.ErrorMessage{"A mealplan exists with this name already "}
  			json.NewEncoder(response).Encode(error)
  		}
  	}
  }
}



func GetMealplans(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func (response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will get all mealplans")
  }
}

func GetMealplanById(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will get mealplan with its id")
  }
}

func GetMealplansByUserId(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will get a specific users mealplans")
  }
}

func AddMealToMealplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will add a meal to a mealplan")
  }
}

func DeleteMealplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will delete a mealplan")
  }
}

func DeleteMealFromMealplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
  	fmt.Println("This will remove a meal from a mealplan")
  }
}
