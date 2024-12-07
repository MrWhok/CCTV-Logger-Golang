package service

import (
	"CCTV-Logger-Golang/src/internal/api/repository"
	"CCTV-Logger-Golang/src/internal/pkg/models"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService interface {
	GetAllImages() ([]models.Image, error)
	GetImage(image *models.Image) map[string]interface{}
	UploadImage(filePath string, totalEntity int) (map[string]string, error)
	DeleteImage(image *models.Image) (map[string]string, error)
	FavoriteImage(image *models.Image, starred bool) (map[string]interface{}, error)
	FindByID(id primitive.ObjectID) (*models.Image, error)
}

type imageService struct {
	repo repository.ImageRepository
}

func NewImageService(repo repository.ImageRepository) ImageService {
	return &imageService{
		repo: repo,
	}
}

func (s *imageService) FindByID(id primitive.ObjectID) (*models.Image, error) {
	return s.repo.FindByID(id)
}

func (s *imageService) GetAllImages() ([]models.Image, error) {
	return s.repo.GetAllImages()
}

func (s *imageService) GetImage(image *models.Image) map[string]interface{} {
	return map[string]interface{}{
		"message":   "Image found successfully",
		"imageData": image,
	}
}

func (s *imageService) UploadImage(filePath string, totalEntity int) (map[string]string, error) {
	fileName := filepath.Base(filePath)
	destPath := filepath.Join("public/staticimages", fileName)

	err := os.Rename(filePath, destPath)
	if err != nil {
		return nil, err
	}

	image := models.NewImage("/staticimages/" + fileName)
	image.TotalEntity = totalEntity

	if err := s.repo.Save(image); err != nil {
		return nil, err
	}

	return map[string]string{
		"message": "File uploaded successfully",
		"url":     "/staticimages/" + fileName,
	}, nil
}

func (s *imageService) DeleteImage(image *models.Image) (map[string]string, error) {
	err := os.Remove(filepath.Join("public", image.ImageURL))
	if err != nil {
		return nil, err
	}

	if err := s.repo.Delete(image); err != nil {
		return nil, err
	}

	return map[string]string{
		"message":  "Image deleted successfully",
		"imageUrl": image.ImageURL,
	}, nil
}

func (s *imageService) FavoriteImage(image *models.Image, starred bool) (map[string]interface{}, error) {
	image.Starred = starred

	filter := bson.M{"_id": image.ID}
	update := bson.M{"$set": bson.M{"starred": image.Starred}}

	if err := s.repo.UpdateOne(filter, update); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":      "Starred status updated successfully",
		"updatedImage": image,
	}, nil
}
