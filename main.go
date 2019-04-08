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
	"github.com/KatieSchmidt/meal_plan/routes"
)

var client *mongo.Client

func main() {
	fmt.Println("Server started on port 5000")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	//try to do the rest of everything and then cancel
	defer cancel()

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/meal", meals.CreateMeal).Methods("POST")
	router.HandleFunc("/meals", meals.GetMeals).Methods("GET")
	router.HandleFunc("/meals/{id}", meals.GetMealById).Methods("GET")
	router.HandleFunc("/meals/{id}/ingredients", meals.AddIngredientToMeal).Methods("PUT")
	router.HandleFunc("/meals/{id}", meals.DeleteMealById).Methods("DELETE")
	router.HandleFunc("/meals/{id}/ingredients/{ingredient_id}", meals.DeleteIngredientFromMeal).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", router))
}
