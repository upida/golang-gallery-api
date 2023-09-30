package controllers

import (
	"gallery/models"
	"gallery/utils/file"
	"gallery/utils/token"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUser(ctx *gin.Context) {

	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

type LoginInput struct {
	Username string `json:"username" binding:"required_without=Email,omitempty"`
	Email    string `json:"email" binding:"required_without=Username,omitempty,email"`
	Password string `json:"password" binding:"required"`
}

func Login(ctx *gin.Context) {

	var input LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	var identifier string
	if input.Email != "" {
		u.Email = input.Email
		identifier = input.Email
	} else {
		u.Username = input.Username
		identifier = input.Username
	}

	u.Password = input.Password

	token, err := models.LoginCheck(identifier, u.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect authentication data."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})

}

type RegisterInput struct {
	Email    string `form:"email" binding:"required,email"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Register(ctx *gin.Context) {

	var input RegisterInput

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u models.User

	u.Email = input.Email
	u.Username = input.Username
	u.Password = input.Password

	user, err := u.SaveUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo_profile, _ := ctx.FormFile("photo_profile")

	if photo_profile != nil {
		var photos []*multipart.FileHeader
		photos = append(photos, photo_profile)

		resultPhotos, err := file.UploadPhotos(ctx, user.ID, photos)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo metadata"})
			return
		}

		user.PhotoProfile = resultPhotos[0].UUID
		user, err = user.UpdateUser()

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registration success", "data": user})

}

type UpdateUserInput struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
}

func UpdateUser(ctx *gin.Context) {
	user_id, err := token.ExtractTokenID(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateUserInput

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByID(user_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Password != "" {
		user.Password = input.Password
		err := user.BeforeSave()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	user, err = user.UpdateUser()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo_profile, _ := ctx.FormFile("photo_profile")

	if photo_profile != nil {
		var photos []*multipart.FileHeader
		photos = append(photos, photo_profile)

		resultPhotos, err := file.UploadPhotos(ctx, user.ID, photos)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo metadata"})
			return
		}
		log.Print(user)
		user.PhotoProfile = resultPhotos[0].UUID
		user, err = user.UpdateUser()

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated success", "data": user})

}
