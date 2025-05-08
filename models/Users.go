package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct represents the User model
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserName       string             `bson:"userName" validate:"required,min=10,max=70"`
	Email          string             `bson:"email" validate:"required,email"`
	Password       string             `bson:"password" validate:"required"`
	Age            int                `bson:"age" validate:"required"`
	UserVerified   VerifiedStatus     `bson:"userVerified"`
	UserVerifyToken VerifyToken       `bson:"userVerifyToken"`
}

// VerifiedStatus struct for user verification flags
type VerifiedStatus struct {
	Email bool `bson:"email"`
	Phone bool `bson:"phone"`
}

// VerifyToken struct for verification tokens
type VerifyToken struct {
	Email *string `bson:"email,omitempty"`
	Phone *string `bson:"phone,omitempty"`
}