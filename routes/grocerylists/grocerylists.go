package grocerylists

import (
	"fmt"
	"encoding/json"
	"context"
  "strings"
  "log"
  "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
  "net/http"
	"github.com/KatieSchmidt/meal_plan/models"
)

func CreateGrocerylist(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
    request.ParseForm()
    response.Header().Set("content-type", "application/x-www-form-urlencoded")

    params := mux.Vars(request)
    mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])

    grocerylist_filter := bson.D{{"associatedmealplanid", mealplan_id}}
    mealplan_filter := bson.D{{"_id", mealplan_id}}

    grocery_collection := mongoClient.Database("go_meals").Collection("grocerylists")
    mealplan_collection := mongoClient.Database("go_meals").Collection("mealplans")

    var mealplan models.Mealplan

    mealplan_err := mealplan_collection.FindOne(ctx, mealplan_filter).Decode(&mealplan)

    //if the mealplan doesnt exist return an error
    if mealplan_err != nil {
      var errors models.Errors
      errors.Grocerylist = "The mealplan wasnt found, couldn't make grocery list"
      json.NewEncoder(response).Encode(errors)
    } else { //otherwise  make a temp grocery list to hold the new one
      var grocerylist models.Grocerylist
      grocerylist.AssociatedMealplanId = mealplan_id
      for _, meal_item := range mealplan.Meals {
        for _, ingredient := range meal_item.Ingredients {
          var inserted = false
          if len(grocerylist.Groceries) == 0{
            var tempitem models.Groceryitem
            tempitem.ID = primitive.NewObjectID()
            tempitem.Ingredient = ingredient.Ingredient
            tempitem.Quantity = ingredient.MeasureUnitQuantity
            tempitem.MeasureUnit = ingredient.MeasureUnit

            grocerylist.Groceries = append(grocerylist.Groceries, tempitem)
          } else {
            for index, list_item := range grocerylist.Groceries {
              if strings.ToLower(ingredient.Ingredient) == strings.ToLower(list_item.Ingredient){
                grocerylist.Groceries[index].Quantity += ingredient.MeasureUnitQuantity
                inserted = true
                break
              }
            }
            if inserted == false {
              var tempitem models.Groceryitem
              tempitem.ID = primitive.NewObjectID()
              tempitem.Ingredient = ingredient.Ingredient
              tempitem.Quantity = ingredient.MeasureUnitQuantity
              tempitem.MeasureUnit = ingredient.MeasureUnit

              grocerylist.Groceries = append(grocerylist.Groceries, tempitem)
            }
          }
        }
      }

      var original_grocerylist models.Grocerylist

      grocerylist_err := grocery_collection.FindOne(ctx, grocerylist_filter).Decode(&original_grocerylist)


      if grocerylist_err != nil {
        fmt.Println("no list found, creating new one")

        _, insertion_err := grocery_collection.InsertOne(ctx, grocerylist)
        if insertion_err != nil {
          var response_message models.Errors
          response_message.Grocerylist = "Couldn't make the new grocery list"
          json.NewEncoder(response).Encode(response_message)
        } else {
          json.NewEncoder(response).Encode(grocerylist)
        }
      } else {
        fmt.Println("list found, replacing")

        var replaced_list models.Grocerylist
        replacement_err := grocery_collection.FindOneAndReplace(ctx, grocerylist_filter, grocerylist).Decode(&replaced_list)
        if replacement_err != nil {
          var response_message models.Errors
          response_message.Grocerylist = "Couldn't replace with the new grocery list"
          json.NewEncoder(response).Encode(response_message)
        } else {
          json.NewEncoder(response).Encode(grocerylist)
        }
      }
    }
  }
}

func GetGrocerylists(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    grocerylist_collection := mongoClient.Database("go_meals").Collection("grocerylists")
    cursor, err := grocerylist_collection.Find(ctx, bson.M{})

    if err != nil {
      log.Fatal(err)
    }

    //create a list of meals of struc models.Meal
  	var grocerylists []models.Grocerylist
  	for cursor.Next(ctx) {
  		var grocerylist models.Grocerylist
  		cursor.Decode(&grocerylist)
  		grocerylists = append(grocerylists, grocerylist)
  	}

  	if len(grocerylists) > 0 {
  		json.NewEncoder(response).Encode(grocerylists)

  	} else {
  		//if there are no meals create a message Struct to send back
			var response_message models.Errors
			response_message.Mealplan = "No grocerylists have been created"
  		json.NewEncoder(response).Encode(response_message)
  	}
  }
}

func GetGrocerylistByMealplan(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("conent-type", "application/x-www-form-urlencoded")
    params := mux.Vars(request)
    mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])
    filter := bson.D{{"associatedmealplanid", mealplan_id}}

    collection := mongoClient.Database("go_meals").Collection("grocerylists")
    var grocerylist models.Grocerylist
    err := collection.FindOne(ctx, filter).Decode(&grocerylist)
    if err != nil {
      var error_message models.Errors
      error_message.Grocerylist = "This mealplan doesnt have a grocery list or the mealpla doesnt exist with this id"
      json.NewEncoder(response).Encode(error_message)
    } else {
      json.NewEncoder(response).Encode(grocerylist)
    }

  }
}

func RemoveItemFromGroceryList(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/x-www-form-urlencoded")
    params := mux.Vars(request)
    mealplan_id, _ := primitive.ObjectIDFromHex(params["mealplan_id"])
    grocery_id, _ := primitive.ObjectIDFromHex(params["grocery_id"])
    collection := mongoClient.Database("go_meals").Collection("grocerylists")
    filter := bson.D{{"associatedmealplanid", mealplan_id}}

    var grocerylist models.Grocerylist
    err := collection.FindOne(ctx, filter).Decode(&grocerylist)

    if err != nil {
      var errors models.Errors
      errors.Grocerylist = "this grocery list wasn't found"
    } else {
      deleted := false
      for index, item := range grocerylist.Groceries {
        if item.ID == grocery_id {
          new_grocerylist_slice := append(grocerylist.Groceries[:index], grocerylist.Groceries[index + 1:]...)
          grocerylist.Groceries = new_grocerylist_slice
          deleted = true
          break
        }
      }
      if deleted != true {
        var error_message models.Errors
        error_message.Grocerylist = "This item didnt get deleted."
        json.NewEncoder(response).Encode(error_message)
      } else {
        var replaced_list models.Grocerylist
        err := collection.FindOneAndReplace(ctx, filter, grocerylist).Decode(&replaced_list)

        if err != nil {
          var error_message models.Errors
          error_message.Grocerylist = "couldnt replace the old meal with a new meal"
          json.NewEncoder(response).Encode(error_message)
        } else {
          json.NewEncoder(response).Encode(grocerylist)
        }
      }
    }


  }
}

func DeleteGroceryList(ctx context.Context, mongoClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {

  }
}
