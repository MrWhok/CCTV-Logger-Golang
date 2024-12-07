package repository

import (
	"CCTV-Logger-Golang/src/db"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(username string) error
}

type userRepository struct {
	collection string
}

func NewUserRepository() UserRepository {
	return &userRepository{
		collection: "users",
	}
}

func (r *userRepository) CreateUser(user *models.User) error {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindUserByUsername(username string) (*models.User, error) {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) DeleteUser(username string) error {
	collection := db.GetCollection(db.Client, r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"username": username})
	return err
}
