package models

import (
	"app_template/src/libs"

	"golang.org/x/crypto/bcrypt"

	"errors"
	"log"
	"time"
)

var User *user

type user struct {
	ID              uint `gorm:"primary_key"`
	Username        string
	Password        []byte
	Email           string
	RememberToken   string
	VerifyToken     string
	ResetToken      string
	Status          int
	EmailVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (thisUser *user) List() []user {
	var data = []user{}
	err := libs.DB.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func (thisUser *user) Info(id uint) (user, error) {
	var data user

	if libs.DB.Where("id  = ? ", id).Find(&data).RecordNotFound() {
		return user{}, errors.New("no existe registro con es id")
	}

	return data, nil
}

func (thisUser *user) Add(username string, password string, email string) (user, error) {
	var data user
	data.Username = username
	data.Email = email

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), 5)

	if hashErr != nil {
		return user{}, hashErr
	}

	data.Password = hash

	if err := libs.DB.Table("users").Create(&data).Error; err != nil {
		return user{}, err
	} else {

		return data, nil
	}
}

func (thisUser *user) Update(id uint) (user, error) {
	var data user

	if libs.DB.Where("id = ? ", id).First(&data).RecordNotFound() {
		return user{}, errors.New("no hay registro con ese id")
	}

	if err := libs.DB.Save(&data).Error; err != nil {
		return user{}, errors.New("no se pudo actualizar")
	}
	return data, nil

}

func (thisUser *user) Del(id uint) error {
	var data user
	if err := libs.DB.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}
