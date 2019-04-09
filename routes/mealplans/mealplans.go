package mealplans

import (
	"fmt"
	// "encoding/json"
	// "context"
	// "time"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
  "net/http"
	// "github.com/KatieSchmidt/meal_plan/models"
)


// var client *mongo.Client

func CreateMealplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all mealplans")
}

// func CreateMealplan(response http.ResponseWriter, request *http.Request) {
// 	request.ParseForm()
// 	response.Header().Set("content-type", "application/x-www-form-urlencoded")
//
// 	if len(request.FormValue("userid")) == 0 || len(request.FormValue("planname")) == 0 {
// 		meal_error := models.ErrorMessage{"One of your form fields was empty"}
// 		json.NewEncoder(response).Encode(meal_error)
// 	} else {
// 		// look for a mealplan that has same name and user
// 			// if it does return an error if not, make the meal
// 		collection := client.Database("go_meals").Collection("meals")
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		var mealplan models.Mealplan
// 		mealplan.UserId = request.FormValue("userid")
// 		mealplan.Planname = request.FormValue("planname")
// 		filter := bson.D{{"userid", mealplan.UserId}, {"planname", mealplan.Planname}}
// 		error_msg := collection.FindOne(ctx, filter)
//
// 		if error_msg != nil {
// 			_, err := collection.InsertOne(ctx, mealplan)
// 			if err != nil {
// 				response_message := models.ErrorMessage{"ERROR: there was an error creating your mealplan"}
// 				json.NewEncoder(response).Encode(response_message)
// 			} else {
// 				//if there isnt an error, meal was inserted, so return the meal
// 				json.NewEncoder(response).Encode(mealplan)
// 			}
// 		} else {
// 			error := models.ErrorMessage{"A mealplan exists with this name already "}
// 			json.NewEncoder(response).Encode(error)
// 		}
// 	}
// }

func GetMealplans(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all mealplans")
}

func GetMealplanById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get mealplan with its id")
}
func GetMealplansByUserId(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get a specific users mealplans")
}


func AddMealToMealplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will add a meal to a mealplan")
}

func DeleteMealplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will delete a mealplan")
}
func DeleteMealFromMealplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will remove a meal from a mealplan")
}
