package models

import (
	"errors"
	"html"
	"strings"

	"gallery/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique" json:"email"`
	Username     string `gorm:"size:255;not null;unique" json:"username"`
	Password     string `gorm:"size:255;not null;" json:"password"`
	PhotoProfile string `gorm:"size:255;" json:"photo_profile"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(identifier string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", identifier).Or("email = ?", identifier).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u User) SaveUser() (User, error) {

	var err error

	err = u.BeforeSave()
	if err != nil {
		return User{}, err
	}

	err = DB.Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPasswordString := string(hashedPassword)
	err = VerifyPassword(u.Password, hashedPasswordString)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	u.Password = hashedPasswordString

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func (u User) UpdateUser() (User, error) {

	var err error

	err = DB.Updates(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}
