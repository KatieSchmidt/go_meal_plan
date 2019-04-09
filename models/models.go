package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Meal struct {
	Username string `json:"username" bson:"username"`
	Mealname  string `json:"mealname" bson:"mealname"`
	TotalCalories int64 `json:"totalcalories" bson:"totalcalories"`
	Ingredients []Ingredient `json:"ingredients" bson:"ingredients"`
	DateAdded time.Time
}

type Ingredient struct {
	ID primitive.ObjectID `json:"ing_id" bson:"ing_id"`
	Ingredient string `json:"ingredient" bson:"ingredient"`
	Calories int64 `json:"calories" bson:"calories"`
	MeasureUnitQuantity int64 `json:"measureunitquantity" bson:"measureunitquantity"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
}

type ErrorMessage struct {
	Error string
}

type User struct {
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Token string `json:"token" bson:"token"`
}

type Mealplan struct {
	UserId string `json:"userid" bson:"userid"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories int64 `json:"totalcalories" bson:"totalcalories"`
	Meals []Meal `json:"meals" bson:"meals"`
}

type Weekplan struct {
	UserId string `json:"user" bson:"user"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories int64 `json:"totalcalories" bson:"totalcalories"`
	Mealplans []Mealplan `json:"mealplans" bson:"mealplans"`
}

type GroceryList struct {
	Mealplan primitive.ObjectID `json:"mealplan" bson:"mealplan"`
	Groceries map[int]GroceryItem `json:"groceries" bson:"groceries"`
}

type WeekList struct {
	Weekplan primitive.ObjectID `json:"weekplan" bson:"weekplan"`
	Groceries map[int]GroceryItem `json:"groceries" bson:"groceries"`
}

type GroceryItem struct {
	Ingredient string `json:"ingredient" bson:"ingredient"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
}
