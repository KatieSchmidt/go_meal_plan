package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/KatieSchmidt/meal_plan/routes/meals"
	"github.com/KatieSchmidt/meal_plan/routes/mealplans"
	"github.com/KatieSchmidt/meal_plan/routes/grocerylists"
	"github.com/KatieSchmidt/meal_plan/routes/users"
	"github.com/KatieSchmidt/meal_plan/routes/weekplans"
)

var client *mongo.Client

func main() {
	fmt.Println("Server started on port 5000")
	ctx := context.Background()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	router := mux.NewRouter()

	//meals routes
	router.HandleFunc("/meals", meals.GetMeals(ctx, client)).Methods("GET")

	router.HandleFunc("/meals/usermeals", meals.GetCurrentUsersMeals(ctx, client)).Methods("GET")

	router.HandleFunc("/meals/{meal_id}", meals.GetMealById(ctx, client)).Methods("GET")

	router.HandleFunc("/meals", meals.CreateMeal(ctx, client)).Methods("POST")

	router.HandleFunc("/meals/{meal_id}", meals.AddIngredientToMeal(ctx, client)).Methods("PUT")

	router.HandleFunc("/meals/{meal_id}/remove/{ingredient_id}", meals.DeleteIngredientFromMeal(ctx, client)).Methods("PUT")

	router.HandleFunc("/meals/{meal_id}", meals.DeleteMealById(ctx, client)).Methods("DELETE")


	//mealplan routes
	router.HandleFunc("/mealplans", mealplans.CreateMealplan(ctx, client)).Methods("POST")

	router.HandleFunc("/mealplans",mealplans.GetMealplans(ctx, client)).Methods("GET")

	router.HandleFunc("/mealplans/{mealplan_id}",mealplans.GetMealplanById(ctx, client)).Methods("GET")

	router.HandleFunc("/mealplans/usermealplans",mealplans.GetCurrentUsersMealplans(ctx, client)).Methods("GET")

	router.HandleFunc("/mealplans/{mealplan_id}/{meal_id}",mealplans.AddMealToMealplan(ctx, client)).Methods("PUT")

	router.HandleFunc("/mealplans/{mealplan_id}",mealplans.DeleteMealplan(ctx, client)).Methods("DELETE")

	router.HandleFunc("/mealplans/{mealplan_id}/remove/{meal_id}",mealplans.DeleteMealFromMealplan(ctx, client)).Methods("PUT")

	//grocerylist routes
	router.HandleFunc("/grocerylists/{mealplan_id}", grocerylists.CreateGrocerylist(ctx, client)).Methods("POST")

	router.HandleFunc("/grocerylists", grocerylists.GetGrocerylists(ctx, client)).Methods("GET")

	router.HandleFunc("/grocerylists/{mealplan_id}", grocerylists.GetGrocerylistByMealplan(ctx, client)).Methods("GET")

	router.HandleFunc("/grocerylists/{mealplan_id}/{grocery_id}", grocerylists.RemoveItemFromGroceryList(ctx, client)).Methods("PUT")

	router.HandleFunc("/grocerylists/{mealplan_id}", grocerylists.DeleteGroceryList(ctx, client)).Methods("DELETE")

	//users routes
	router.HandleFunc("/users/register", users.RegisterUser(ctx, client)).Methods("POST")

	router.HandleFunc("/users/login", users.LoginUser(ctx, client)).Methods("POST")

	//weekplan routes
	router.HandleFunc("/weekplans", weekplans.CreateWeekplan(ctx, client)).Methods("POST")

	router.HandleFunc("/weekplans", weekplans.GetWeekplans(ctx, client)).Methods("GET")

	router.HandleFunc("/weekplans/userweekplans", weekplans.GetCurrentUsersWeekplans(ctx, client)).Methods("GET")

	router.HandleFunc("/weekplans/{weekplan_id}", weekplans.GetWeekplanById(ctx, client)).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", router))

}
