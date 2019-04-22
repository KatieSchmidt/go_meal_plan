package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/dgrijalva/jwt-go"
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
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Mealplan string `json:"mealplan,omitempty" bson:"mealplan,omitempty"`
	Weekplan string `json:"weekplan,omitempty" bson:"weekplan,omitempty"`
	Meal string `json:"meal,omitempty" bson:"meal,omitempty"`
	Ingredient string `json:"ingredient,omitempty" bson:"ingredient,omitempty"`
	Grocerylist string `json:"grocerylist,omitempty" bson:"grocerylist,omitempty"`
}

type Mealplan struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	User primitive.ObjectID `json:"user" bson:"user"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories float64 `json:"totalcalories" bson:"totalcalories"`
	Meals []Meal `json:"meals" bson:"meals"`
}

type Weekplan struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	User primitive.ObjectID `json:"user" bson:"user"`
	Planname string `json:"planname" bson:"planname"`
	TotalCalories float64 `json:"totalcalories" bson:"totalcalories"`
	Mealplans []Mealplan `json:"mealplans" bson:"mealplans"`
}

type Grocerylist struct {
	AssociatedMealplanId primitive.ObjectID  `json:"associatedmealplanid" bson:"associatedmealplanid"`
	Groceries []Groceryitem
}

type Groceryitem struct {
	ID primitive.ObjectID  `json:"_id" bson:"_id"`
	Ingredient string `json:"ingredient" bson:"ingredient"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	MeasureUnit string `json:"measureunit" bson:"measureunit"`
}

type User struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password []byte `json:"password" bson:"password"`
}

type Claims struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

type JOT struct {
	Success bool `json:"success", bson:"success"`
	Token string `json:"token", bson:"token"`
}
