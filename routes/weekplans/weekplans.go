package weekplans

import (
	"fmt"
  "net/http"
)

func GetWeekplans(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get all weekplans")
}

func GetWeekplanById(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get weekplan with its id")
}
func GetWeekplansByUserId(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will get a specific users weekplans")
}

func CreateWeekplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will create a new weekplan")
}
func AddMealplanToWeekplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will add a mealplan to a weekplan")
}

func DeleteWeekplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will delete a weekplan")
}
func DeleteMealplanFromWeekplan(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will remove a mealplan from a weekplan")
}
