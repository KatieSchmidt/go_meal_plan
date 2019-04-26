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

//done
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

//done
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

//done
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

//done
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

//done
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

func AddMealplanToWeekplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		// response.Header().Set("content-type", "application/x-www-form-urlencoded")
		headerTkn := request.Header.Get("Authorization")
		var newClaims models.Claims
		token, _ := jwt.ParseWithClaims(headerTkn, &newClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte("my_secret_key"), nil //will be hidden in production
    })

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			//////////
	  	params := mux.Vars(request)

	  	mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])
	    weekplan_id, _ := primitive.ObjectIDFromHex(params["weekplan_id"])

	  	mealplanfilter := bson.D{{"user", claims.ID},{"_id", mealplan_id }}
	    weekplanfilter := bson.D{{"user", claims.ID},{"_id", weekplan_id }}

	  	mealplancollection := mongoClient.Database("go_meals").Collection("mealplans")
	    weekplancollection := mongoClient.Database("go_meals").Collection("weekplans")

	  	var resulting_mealplan models.Mealplan
	    var resulting_weekplan models.Weekplan

	  	mealplan_error_msg := mealplancollection.FindOne(ctx, mealplanfilter).Decode(&resulting_mealplan)
	  	if mealplan_error_msg != nil {
				var response_message models.Errors
				response_message.Mealplan = "Mealplan not found"
	  		json.NewEncoder(response).Encode(response_message)
	  	} else {
	      weekplan_error_msg := weekplancollection.FindOne(ctx, weekplanfilter).Decode(&resulting_weekplan)

	  		if weekplan_error_msg != nil  {
					var response_message models.Errors
					response_message.Mealplan = "You cant add a mealplan to a non-existant weekplan"
		  		json.NewEncoder(response).Encode(response_message)
	  		} else {
	  			resulting_weekplan.Mealplans = append(resulting_weekplan.Mealplans, resulting_mealplan)
	        resulting_weekplan.TotalCalories = resulting_weekplan.TotalCalories + resulting_mealplan.TotalCalories

	        var updated_weekplan models.Weekplan
	        insertion_error := weekplancollection.FindOneAndReplace(ctx, weekplanfilter, resulting_weekplan).Decode(&updated_weekplan)

	        if insertion_error != nil {
						var response_message models.Errors
						response_message.Weekplan = "Unable to add mealplan to weekplan"
			  		json.NewEncoder(response).Encode(response_message)
	        } else {
	          err := weekplancollection.FindOne(ctx, weekplanfilter).Decode(&updated_weekplan)
	          if err != nil {
							var response_message models.Errors
							response_message.Weekplan = "Couldnt find updated weekplan"
				  		json.NewEncoder(response).Encode(response_message)
	          } else {
	            json.NewEncoder(response).Encode(updated_weekplan)
	          }
	        }
	  		}
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

		request.ParseForm()
		// response.Header().Set("content-type", "application/x-www-form-urlencoded")
		headerTkn := request.Header.Get("Authorization")
		var newClaims models.Claims
		token, _ := jwt.ParseWithClaims(headerTkn, &newClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte("my_secret_key"), nil //will be hidden in production
		})

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {

			params := mux.Vars(request)

			mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])
			weekplan_id, _ := primitive.ObjectIDFromHex(params["weekplan_id"])

			mealplanfilter := bson.D{{"user", claims.ID},{"_id", mealplan_id }}
			weekplanfilter := bson.D{{"user", claims.ID},{"_id", weekplan_id }}

			mealplancollection := mongoClient.Database("go_meals").Collection("mealplans")
			weekplancollection := mongoClient.Database("go_meals").Collection("weekplans")

			var resulting_mealplan models.Mealplan
			var resulting_weekplan models.Weekplan

			mealplan_error_msg := mealplancollection.FindOne(ctx, mealplanfilter).Decode(&resulting_mealplan)
			if mealplan_error_msg != nil {
				var response_message models.Errors
				response_message.Mealplan = "Mealplan not found"
				json.NewEncoder(response).Encode(response_message)
			} else {
				weekplan_error_msg := weekplancollection.FindOne(ctx, weekplanfilter).Decode(&resulting_weekplan)

				if weekplan_error_msg != nil  {
					var response_message models.Errors
					response_message.Mealplan = "You cant add a mealplan to a non-existant weekplan"
					json.NewEncoder(response).Encode(response_message)
				} else {
					inserted := false
					for index, mealplan := range resulting_weekplan.Mealplans {
						if mealplan.ID == resulting_mealplan.ID {
							new_mealplan_slice := append(resulting_weekplan.Mealplans[:index], resulting_weekplan.Mealplans[index + 1:]...)
							resulting_weekplan.Mealplans = new_mealplan_slice
							resulting_weekplan.TotalCalories = resulting_weekplan.TotalCalories - resulting_mealplan.TotalCalories
							inserted = true
							break
						}
					}
					if inserted != true {
						var error_msg models.Errors
						error_msg.Weekplan = "this mealplan wasn't in the weekplan"
						json.NewEncoder(response).Encode(error_msg)
					} else {
						var updated_weekplan models.Weekplan
						insertion_error := weekplancollection.FindOneAndReplace(ctx, weekplanfilter, resulting_weekplan).Decode(&updated_weekplan)

						if insertion_error != nil {
							var response_message models.Errors
							response_message.Weekplan = "Unable to remove mealplan from weekplan"
							json.NewEncoder(response).Encode(response_message)
						} else {
							err := weekplancollection.FindOne(ctx, weekplanfilter).Decode(&updated_weekplan)
							if err != nil {
								var response_message models.Errors
								response_message.Weekplan = "Couldnt find updated weekplan"
								json.NewEncoder(response).Encode(response_message)
							} else {
								json.NewEncoder(response).Encode(updated_weekplan)
							}
						}
					}
				}
			}
		} else {
			var response_message models.Errors
			response_message.Weekplan = "there was an error"
			json.NewEncoder(response).Encode(response_message)
		}
  }
}
