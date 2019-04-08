package users

import (
	"fmt"
  "net/http"
)

func RegisterUser(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will register a user")
}

func LoginUser(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This will login a user")
}
