package weekplans

import (
	// "fmt"
	"encoding/json"
	"context"
  "log"
  "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "net/http"
	"github.com/KatieSchmidt/meal_plan/models"
	"github.com/dgrijalva/jwt-go"
)

func  CreateWeekplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		headerToken := request.Header.Get("Authorization")
  	request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")
		var newClaims models.Claims
		token, _ := jwt.ParseWithClaims(headerToken, &newClaims, func(token *jwt.Token)(interface{}, error){
			return []byte("my_secret_key"), nil
		})
		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			if len(request.FormValue("planname")) == 0 {
				var plan_errors models.Errors
				plan_errors.Planname = "Planname required"
	  		json.NewEncoder(response).Encode(plan_errors)
	  	} else {
	  		// look for a mealplan that has same name and user
	  			// if it does return an error if not, make the meal
	  		collection := mongoClient.Database("go_meals").Collection("weekplans")
	  		var weekplan models.Weekplan
				weekplan.ID = primitive.NewObjectID()
	  		weekplan.User = claims.ID
	  		weekplan.Planname = request.FormValue("planname")
	  		filter := bson.D{{"user", weekplan.User}, {"planname", weekplan.Planname}}

				var temp models.Weekplan
	  		error_msg := collection.FindOne(ctx, filter).Decode(&temp)

	  		if error_msg != nil {
					_, err := collection.InsertOne(ctx, weekplan)
	  			if err != nil {

						var response_message models.Errors
						response_message.Mealplan = "Could not create weekplan"

	  				json.NewEncoder(response).Encode(response_message)
	  			} else {
	  				//if there isnt an error, meal was inserted, so return the meal
	  				json.NewEncoder(response).Encode(weekplan)
	  			}
	  		} else {
					var response_message models.Errors
					response_message.Weekplan = "A weekplan exists with this name already "
	  			json.NewEncoder(response).Encode(response_message)
	  		}
	  	}
		} else {
			response_message := models.ErrorMessage{"Weekplan couldnt be created"}
			json.NewEncoder(response).Encode(response_message)
		}
  }
}



func GetWeekplans(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func (response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
  	collection := mongoClient.Database("go_meals").Collection("weekplans")

  	cursor, err := collection.Find(ctx, bson.M{})

  	if err != nil {
  		log.Fatal(err)
  	}

  	//create a list of meals of struc models.Meal
  	var weekplans []models.Weekplan
  	for cursor.Next(ctx) {
  		var weekplan models.Weekplan
  		cursor.Decode(&weekplan)
  		weekplans = append(weekplans, weekplan)
  	}

  	if len(weekplans) > 0 {
  		json.NewEncoder(response).Encode(weekplans)

  	} else {
  		//if there are no meals create a message Struct to send back
			var response_message models.Errors
			response_message.Weekplan = "No weekplans have been created"
  		json.NewEncoder(response).Encode(response_message)
  	}
  }
}

func GetWeekplanById(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		headerToken := request.Header.Get("Authorization")
  	request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")
		var newClaims models.Claims
		token, _ := jwt.ParseWithClaims(headerToken, &newClaims, func(token *jwt.Token)(interface{}, error){
			return []byte("my_secret_key"), nil
		})
		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			params := mux.Vars(request)
			collection := mongoClient.Database("go_meals").Collection("weekplans")

			var resulting_weekplan models.Weekplan
			id, _ := primitive.ObjectIDFromHex(params["weekplan_id"])
			filter := bson.D{{"user", claims.ID}, {"_id", id}}
			error_msg := collection.FindOne(ctx, filter).Decode(&resulting_weekplan)

			if error_msg != nil {
				var response_message models.Errors
				response_message.Weekplan = "Weekplan with this id doesn't exist for this user"
				json.NewEncoder(response).Encode(response_message)
			} else {
				json.NewEncoder(response).Encode(resulting_weekplan)
			}
		} else {
			var response_message models.Errors
			response_message.Weekplan = "You arent logged in"
  		json.NewEncoder(response).Encode(response_message)
		}

  }
}

func AddMealplanToWeekplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
    request.ParseForm()
  	response.Header().Set("content-type", "application/x-www-form-urlencoded")

  	params := mux.Vars(request)

  	meal_id, _ := primitive.ObjectIDFromHex(params["meal_id"])
    plan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])

  	mealfilter := bson.D{{"_id", meal_id }}
    mealplanfilter := bson.D{{"_id", plan_id }}

  	mealcollection := mongoClient.Database("go_meals").Collection("meals")
    mealplancollection := mongoClient.Database("go_meals").Collection("mealplans")

  	var resulting_meal models.Meal
    var resulting_mealplan models.Mealplan

  	meal_error_msg := mealcollection.FindOne(ctx, mealfilter).Decode(&resulting_meal)
  	if meal_error_msg != nil {
			var response_message models.Errors
			response_message.Meal = "Meal not found"
  		json.NewEncoder(response).Encode(response_message)
  	} else {
      mealplan_error_msg := mealplancollection.FindOne(ctx, mealplanfilter).Decode(&resulting_mealplan)

  		if mealplan_error_msg != nil  {
				var response_message models.Errors
				response_message.Mealplan = "You cant add a meal to a non-existant mealplan"
	  		json.NewEncoder(response).Encode(response_message)
  		} else {
  			resulting_mealplan.Meals = append(resulting_mealplan.Meals, resulting_meal)
        resulting_mealplan.TotalCalories = resulting_mealplan.TotalCalories + resulting_meal.TotalCalories

        var updated_mealplan models.Mealplan
        insertion_error := mealplancollection.FindOneAndReplace(ctx, mealplanfilter, resulting_mealplan).Decode(&updated_mealplan)

        if insertion_error != nil {
					var response_message models.Errors
					response_message.Mealplan = "Unable to add meal to mealplan"
		  		json.NewEncoder(response).Encode(response_message)
        } else {
          err := mealplancollection.FindOne(ctx, mealplanfilter).Decode(&updated_mealplan)
          if err != nil {
						var response_message models.Errors
						response_message.Mealplan = "Couldnt find updated mealplan"
			  		json.NewEncoder(response).Encode(response_message)
          } else {
            json.NewEncoder(response).Encode(updated_mealplan)
          }
        }
  		}
  	}
  }
}

func GetCurrentUsersWeekplans(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		headerTkn := request.Header.Get("Authorization")
		var newClaims models.Claims
		token, _ := jwt.ParseWithClaims(headerTkn, &newClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte("my_secret_key"), nil //will be hidden in production
    })

    if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			collection := mongoClient.Database("go_meals").Collection("weekplans")

	    filter := bson.D{{"user", claims.ID}}

	    cursor, err := collection.Find(ctx, filter)

	    if err != nil {
				var response_message models.Errors
				response_message.Weekplan = "Weekplans not found for this user"
				json.NewEncoder(response).Encode(response_message)
	    } else {
	      var weekplans []models.Weekplan
	      for cursor.Next(ctx){
	        var weekplan models.Weekplan
	        cursor.Decode(&weekplan)
	        weekplans = append(weekplans, weekplan)
	      }
	      if len(weekplans) > 0 {
	        json.NewEncoder(response).Encode(weekplans)
	      } else {
					var response_message models.Errors
					response_message.Weekplan = "This user has no weekplans"
					json.NewEncoder(response).Encode(response_message)
	      }
	    }
    } else {
      var response_message models.Errors
			response_message.Weekplan = "there was an error"
			json.NewEncoder(response).Encode(response_message)
    }
  }
}

func DeleteWeekplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		headerTkn := request.Header.Get("Authorization")
		var newClaims models.Claims

		token, _ := jwt.ParseWithClaims(headerTkn, &newClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte("my_secret_key"), nil //will be hidden in production
    })

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			params := mux.Vars(request)
	    collection := mongoClient.Database("go_meals").Collection("weekplans")
			var user_id = claims.ID
	    id, _ := primitive.ObjectIDFromHex(params["weekplan_id"])
	    filter := bson.D{{"_id", id }, {"user", user_id}}
	    result, _ := collection.DeleteOne(ctx, filter)

	    // DeleteOne always returns a result. error is always nil. so check to see if deleted count is equal to zero instead
	    if result.DeletedCount == 0 {
				var response_message models.Errors
				response_message.Mealplan = "weekplan not found"
				json.NewEncoder(response).Encode(response_message)
	    } else {
	      response_message := models.ResponseMessage{"weekplan deleted"}
	      json.NewEncoder(response).Encode(response_message)
	    }

		} else {
			var response_message models.Errors
			response_message.Weekplan = "there was an error"
			json.NewEncoder(response).Encode(response_message)
		}
  }
}

func DeleteMealplanFromWeekplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {

    params := mux.Vars(request)

  	mealplan_collection := mongoClient.Database("go_meals").Collection("mealplans")
    meal_collection := mongoClient.Database("go_meals").Collection("meals")

    mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])
    meal_id, _ := primitive.ObjectIDFromHex(params["meal_id"])

    mealplan_filter := bson.D{{"_id", mealplan_id}}
    meal_filter := bson.D{{"_id", meal_id}}

    var original_mealplan models.Mealplan
    err := mealplan_collection.FindOne(ctx, mealplan_filter).Decode(&original_mealplan)
    if err != nil {
			var response_message models.Errors
			response_message.Mealplan = "Mealplan Not Found"
			json.NewEncoder(response).Encode(response_message)
    } else {
      var meal models.Meal
      err := meal_collection.FindOne(ctx, meal_filter).Decode(&meal)

      if err != nil {
				var response_message models.Errors
				response_message.Mealplan = "Meal Not Found"
				json.NewEncoder(response).Encode(response_message)
      } else {
        //even if there is more than one occurance of the meal, only delete one of them from the plan at a time
        index := 0
        running := true
        var calories_to_remove float64 = 0
        for i := 0; i < len(original_mealplan.Meals) && running == true; i ++ {
          if original_mealplan.Meals[i].ID == meal.ID {
            index = i
            calories_to_remove = original_mealplan.Meals[i].TotalCalories
            running = false
          }
      }
      if running == true && calories_to_remove == 0{
				var response_message models.Errors
				response_message.Mealplan = "Meal Not Found in mealplan"
				json.NewEncoder(response).Encode(response_message)
      } else {
        new_meals_slice := append(original_mealplan.Meals[:index], original_mealplan.Meals[index + 1:]...)

        original_mealplan.Meals = new_meals_slice
        original_mealplan.TotalCalories = original_mealplan.TotalCalories - calories_to_remove

        var updated_mealplan models.Mealplan
        error_msg := mealplan_collection.FindOneAndReplace(ctx, mealplan_filter, original_mealplan).Decode(&updated_mealplan)
        if error_msg != nil  {
					var response_message models.Errors
					response_message.Mealplan = "Unable to remove meal from mealplan"
					json.NewEncoder(response).Encode(response_message)
        } else {
          json.NewEncoder(response).Encode(original_mealplan)
        }
      }
      }
    }
  }
}
