package file

import (
	"errors"
	"gallery/models"
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadPhotos(ctx *gin.Context, user_id uint, photos []*multipart.FileHeader) ([]*models.Photo, error) {
	var photosModel []*models.Photo

	for _, photo := range photos {
		uniqid := uuid.New().String()
		filename := uniqid + "_" + photo.Filename
		filePath := filepath.Join("uploads", filename)
		if err := ctx.SaveUploadedFile(photo, filePath); err != nil {
			return nil, errors.New("Failed to save photo!")
		}
		photosModel = append(photosModel, &models.Photo{
			UserID:   user_id,
			Title:    photo.Filename,
			Caption:  "",
			UUID:     uniqid,
			Filename: filename,
		})
	}

	err := models.InsertPhotos(photosModel)
	return photosModel, err
}
