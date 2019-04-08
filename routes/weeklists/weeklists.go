package weeklists

import (
	"fmt"
  "net/http"
)

func GetWeeklists(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all weekly grocery lists")
}

func GetWeeklistById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get weekly grocery list with its id")
}
func GetWeeklistsByUserId(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get a specific users weekly grocery lists")
}

func CreateWeeklist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will create a new weekly grocery list")
}
func AddWeekplanToWeeklist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will add a weekly mealplan to a weekly grocery list")
}

func DeleteWeeklist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will delete a weekly grocery list")
}
func DeleteWeekplanFromWeeklist(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will remove a weekly mealplan from a weekly grocery list")
}
