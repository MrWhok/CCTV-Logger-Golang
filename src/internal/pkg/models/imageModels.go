package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image represents an image document in MongoDB
type Image struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ImageURL    string             `bson:"imageUrl,omitempty"`
	TotalEntity int                `bson:"totalEntity,omitempty"`
	Date        time.Time          `bson:"date,omitempty"`
	Time        string             `bson:"time,omitempty"`
	Starred     bool               `bson:"starred,omitempty"`
}

// NewImage creates a new Image instance with default values
func NewImage(imageURL string) *Image {
	return &Image{
		ID:          primitive.NewObjectID(),
		ImageURL:    imageURL,
		TotalEntity: 0,
		Date:        time.Now(),
		Time:        time.Now().Format("15:04:05"),
		Starred:     false,
	}
}
