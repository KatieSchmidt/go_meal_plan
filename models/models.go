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
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Ingredient string `json:"ingredient" bson:"ingredient"`
	Calories float64 `json:"calories" bson:"calories"`
	MeasureUnitQuantity float64 `json:"measureunitquantity" bson:"measureunitquantity"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
}

type ResponseMessage struct {
	Response string
}
type ErrorMessage struct {
	Response string
}

type Errors struct{
	Planname string `json:"planname,omitempty" bson:"planname,omitempty"`
	User string `json:"user,omitempty" bson:"user,omitempty"`
	Mealplan string `json:"mealplan,omitempty" bson:"mealplan,omitempty"`
	Meal string `json:"meal,omitempty" bson:"meal,omitempty"`
	Ingredient string `json:"ingredient,omitempty" bson:"ingredient,omitempty"`
}

type Mealplan struct {
	User string `json:"user" bson:"user"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories float64 `json:"totalcalories" bson:"totalcalories"`
	Meals []Meal `json:"meals" bson:"meals"`
}
