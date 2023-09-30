package models

import (
	"errors"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Title    string `gorm:"size:255;not null" json:"title"`
	Caption  string `gorm:"text" json:"caption"`
	Filename string `gorm:"not null"`
	UUID     string `gorm:"unique;not null"`
	UserID   uint   `json:"user_id"`
	// Relation
	User User `gorm:"foreignkey:UserID" json:"user"`
}

func GetPhotosByUserID(user_id uint) ([]Photo, error) {
	var photos []Photo

	DB.Where("user_id = ?", user_id).Find(&photos)

	return photos, nil
}

func GetPhotoByUUID(uuid string) (Photo, error) {
	var photo Photo

	err := DB.Where("UUID = ?", uuid).First(&photo).Error
	if err != nil {
		return photo, errors.New("Photo not found!")
	}

	return photo, nil
}

func InsertPhotos(photos []*Photo) error {
	err := DB.Create(photos).Error
	if err != nil {
		return errors.New(err.Error())
	}

	return err
}

func DeletePhotoByUserID(user_id uint, photo Photo) error {
	// Define the path of the file to be deleted
	filePath := filepath.Join("uploads", photo.Filename)
	// Delete the file metadata from the database
	err := DB.Delete(&photo).Error
	if err != nil {
		return errors.New("Failed to delete photo from database!")
	}
	// Delete the file from the server
	err = os.Remove(filePath)
	if err != nil {
		return errors.New("Failed to delete photo from upload folder!")
	}

	return nil
}
