package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName,omitempty" json:"firstName"`
	LastName  string             `bson:"lastName,omitempty" json:"lastName"`
	Username  string             `bson:"username,omitempty" json:"username"`
	Password  string             `bson:"password,omitempty" json:"password"`
}

// UserData represents the user data received in the request
type UserData struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}
