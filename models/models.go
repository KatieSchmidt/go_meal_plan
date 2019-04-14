package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User primitive.ObjectID  `json:"user" bson:"user"`
	Mealname  string `json:"mealname" bson:"mealname"`
	TotalCalories float64 `json:"totalcalories" bson:"totalcalories"`
	Ingredients []Ingredient `json:"ingredients" bson:"ingredients"`
	DateAdded time.Time
}

type Ingredient struct {
	ID primitive.ObjectID `json:"ing_id" bson:"ing_id"`
	Ingredient string `json:"ingredient" bson:"ingredient"`
	Calories float64 `json:"calories" bson:"calories"`
	MeasureUnitQuantity float64 `json:"measureunitquantity" bson:"measureunitquantity"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
}

type ErrorMessage struct {
	Error string
}

type Mealplan struct {
	User string `json:"user" bson:"user"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories float64 `json:"totalcalories" bson:"totalcalories"`
	Meals []Meal `json:"meals" bson:"meals"`
}
