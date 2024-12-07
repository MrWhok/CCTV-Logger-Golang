package repository

import (
	"CCTV-Logger-Golang/src/db"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageRepository interface {
	GetAllImages() ([]models.Image, error)
	Save(image *models.Image) error
	Delete(image *models.Image) error
	FindByID(id primitive.ObjectID) (*models.Image, error)
	UpdateOne(filter interface{}, update interface{}) error
}

type imageRepository struct {
	collection string
}

func NewImageRepository() ImageRepository {
	return &imageRepository{
		collection: "images",
	}
}

func (r *imageRepository) GetAllImages() ([]models.Image, error) {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var images []models.Image
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var image models.Image
		if err := cursor.Decode(&image); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func (r *imageRepository) Save(image *models.Image) error {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, image)
	return err
}

func (r *imageRepository) Delete(image *models.Image) error {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": image.ID})
	return err
}

func (r *imageRepository) FindByID(id primitive.ObjectID) (*models.Image, error) {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var image models.Image
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&image)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *imageRepository) UpdateOne(filter interface{}, update interface{}) error {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
