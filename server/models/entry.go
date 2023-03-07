package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entry struct {
	ID            primitive.ObjectID `bson:"_id"`
	Dish          *string            `json:"dish"`
	Carbohydrates *float64           `json:"carbohydrates"`
	Protein       *float64           `json:"protein"`
	Fat           *float64           `json:"fat"`
	Ingredients   *string            `json:"ingredients"`
	Calories      *string            `json:"calories"`
}
