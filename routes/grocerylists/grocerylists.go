package grocerylists

import (
	"fmt"
  "net/http"
)

func GetGrocerylists(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all grocerylists")
}

func GetGrocerylistById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get grocerylist with its id")
}
func GetGrocerylistsByUserId(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get a specific users grocerylists")
}

func CreateGrocerylist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will create a new grocerylist")
}
func AddMealplanToGrocerylist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will add a mealplan to a grocerylist")
}

func DeleteGrocerylist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will delete a grocerylist")
}
func DeleteMealplanFromGrocerylist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will remove a mealplan from a grocerylist")
}
