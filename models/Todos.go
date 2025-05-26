package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ Todo Model (MongoDB field names properly matched)
type Todo struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Date            string             `bson:"date" validate:"required"`
    TodoNo          int                `bson:"todoNo" validate:"required"`
    TaskTitle       string             `bson:"taskTitle" validate:"required"`
    TaskDescription string             `bson:"taskDescription" validate:"required"`
    Image           string             `bson:"image,omitempty" json:"image,omitempty"` // ✅ Allows optional image upload
}