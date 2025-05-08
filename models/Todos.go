package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Todo struct represents the Todo model
type Todo struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Date           string             `bson:"date" validate:"required"`
	TodoNo         int                `bson:"todoNo" validate:"required"`
	TodoTitle      string             `bson:"todoTitle" validate:"required"`
	TodoDescription string            `bson:"todoDescription" validate:"required"`
	FileUpload     string             `bson:"fileUpload"`
}