package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/KatieSchmidt/meal_plan/routes/meals"
	"github.com/KatieSchmidt/meal_plan/routes/users"
	"github.com/KatieSchmidt/meal_plan/routes/mealplans"
	"github.com/KatieSchmidt/meal_plan/routes/weekplans"
	"github.com/KatieSchmidt/meal_plan/routes/grocerylists"
	"github.com/KatieSchmidt/meal_plan/routes/weeklists"
)

var client *mongo.Client

func main() {
	fmt.Println("Server started on port 5000")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//try to do the rest of everything and then cancel
	defer cancel()

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	//meals routes
	router.HandleFunc("/meal", meals.CreateMeal).Methods("POST")
	router.HandleFunc("/meals", meals.GetMeals).Methods("GET")
	router.HandleFunc("/meals/{id}", meals.GetMealById).Methods("GET")
	router.HandleFunc("/meals/{id}/ingredients", meals.AddIngredientToMeal).Methods("PUT")
	router.HandleFunc("/meals/{id}", meals.DeleteMealById).Methods("DELETE")
	router.HandleFunc("/meals/{id}/ingredients/{ingredient_id}", meals.DeleteIngredientFromMeal).Methods("PUT")
	//user routes
	router.HandleFunc("/users/login", users.LoginUser).Methods("POST")
	router.HandleFunc("/users/register", users.RegisterUser).Methods("GET")
	//mealplan routes
	router.HandleFunc("/mealplans", mealplans.CreateMealplan).Methods("POST")
	router.HandleFunc("/mealplans", mealplans.GetMealplans).Methods("GET")
	router.HandleFunc("/mealplans/{mealplan_id}", mealplans.GetMealplanById).Methods("GET")
	router.HandleFunc("/mealplans/{mealplan_id}/{meal_id}", mealplans.AddMealToMealplan).Methods("PUT")
	router.HandleFunc("/mealplans/{mealplan_id}", mealplans.DeleteMealplan).Methods("DELETE")
	router.HandleFunc("/mealplans/{mealplan_id}/{meal_id}", mealplans.DeleteMealFromMealplan).Methods("PUT")
	//weekplan routes
	router.HandleFunc("/weekplans", weekplans.CreateWeekplan).Methods("POST")
	router.HandleFunc("/weekplans", weekplans.GetWeekplans).Methods("GET")
	router.HandleFunc("/weekplans/{weekplan_id}", weekplans.GetWeekplanById).Methods("GET")
	router.HandleFunc("/weekplans/{weekplan_id}/{mealplan_id}", weekplans.AddMealplanToWeekplan).Methods("PUT")
	router.HandleFunc("/weekplans/{weekplan_id}", weekplans.DeleteWeekplan).Methods("DELETE")
	router.HandleFunc("/weekplans/{weekplan_id}/{mealplan_id}", weekplans.DeleteMealplanFromWeekplan).Methods("PUT")
	//grocery list routes
	router.HandleFunc("/grocerylists", grocerylists.CreateGrocerylist).Methods("POST")
	router.HandleFunc("/grocerylists", grocerylists.GetGrocerylists).Methods("GET")
	router.HandleFunc("/grocerylists/{grocerylist_id}", grocerylists.GetGrocerylistById).Methods("GET")
	router.HandleFunc("/grocerylists/{grocerylist_id}/{mealplan_id}", grocerylists.AddMealplanToGrocerylist).Methods("PUT")
	router.HandleFunc("/grocerylists/{grocerylist_id}", grocerylists.DeleteGrocerylist).Methods("DELETE")
	router.HandleFunc("/grocerylists/{grocerylist_id}/{mealplan_id}", grocerylists.DeleteMealplanFromGrocerylist).Methods("PUT")
	//weekly grocery list routes
	router.HandleFunc("/weeklists", weeklists.CreateWeeklist).Methods("POST")
	router.HandleFunc("/weeklists", weeklists.GetWeeklists).Methods("GET")
	router.HandleFunc("/weeklists/{weeklist_id}", weeklists.GetWeeklistById).Methods("GET")
	router.HandleFunc("/weeklists/{weeklist_id}/{weekplan_id}", weeklists.AddWeekplanToWeeklist).Methods("PUT")
	router.HandleFunc("/weeklists/{weeklist_id}", weeklists.DeleteWeeklist).Methods("DELETE")
	router.HandleFunc("/weeklists/{weeklist_id}/{weekplan_id}", weeklists.DeleteWeekplanFromWeeklist).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", router))
}
