package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	Username string  `json:"username" bson:"username"`
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
