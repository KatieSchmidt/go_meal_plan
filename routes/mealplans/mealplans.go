package mealplans

import (
	"fmt"
  "net/http"
)

func GetMealplans(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all mealplans")
}

func GetMealplanById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get mealplan with its id")
}
func GetMealplansByUserId(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get a specific users mealplans")
}

func CreateMealplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will create a new mealplan")
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
