package controllers

import (
	"fmt"
	"gallery/models"
	"gallery/utils/file"
	"gallery/utils/token"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UserPhotos(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photos, err := models.GetPhotosByUserID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": photos})
	return
}

func UploadPhotos(ctx *gin.Context) {
	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	photos := form.File["photos[]"]

	resultPhotos, err := file.UploadPhotos(ctx, user_id, photos)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo metadata"})
		return
	}
	// Return a success message and the file metadata
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "photos": resultPhotos})
}

// GetFile is a function that retrieves a file from the server
func GetPhoto(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	photo, err := models.GetPhotoByUUID(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	photoPath := filepath.Join("uploads", photo.Filename)
	photoData, err := os.Open(photoPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open photo"})
		return
	}
	defer photoData.Close()

	// Read the first 512 bytes of the file to determine its content type
	photoHeader := make([]byte, 512)
	_, err = photoData.Read(photoHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read photo"})
		return
	}
	photoContentType := http.DetectContentType(photoHeader)

	// Get the file info
	photoInfo, err := photoData.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Set the headers for the file transfer and return the file
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", photo.Title))
	ctx.Header("Content-Type", photoContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", photoInfo.Size()))
	ctx.File(photoPath)
}

// DeleteFile is a function that deletes a file from the server and its metadata from the database
func DeletePhoto(ctx *gin.Context) {
	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get the unique identifier of the file to be deleted
	uuid := ctx.Param("uuid")
	var photo models.Photo
	// Retrieve the file metadata from the database
	photo, err = models.GetPhotoByUUID(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if photo.UserID != user_id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Account is not allowed to delete photo!"})
		return
	}

	err = models.DeletePhotoByUserID(user_id, photo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return a success message
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File " + photo.Filename + " deleted successfully",
	})
}
