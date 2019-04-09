package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/KatieSchmidt/meal_plan/routes"
)

var client *mongo.Client

func main() {
	fmt.Println("Server started on port 5000")
	ctx := context.Background()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	router := mux.NewRouter()

	router.HandleFunc("/meals", meals.CreateMeal(ctx, client)).Methods("POST")
	router.HandleFunc("/meals", meals.GetMeals(ctx, client)).Methods("GET")
	router.HandleFunc("/meals/{id}", meals.GetMealById(ctx, client)).Methods("GET")
	router.HandleFunc("/meals/{id}/ingredients", meals.AddIngredientToMeal(ctx, client)).Methods("PUT")
	router.HandleFunc("/meals/{id}", meals.DeleteMealById(ctx, client)).Methods("DELETE")
	router.HandleFunc("/meals/{id}/ingredients/{ingredient_id}", meals.DeleteIngredientFromMeal(ctx, client)).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", router))
}
