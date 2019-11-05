package models

import "github.com/jinzhu/gorm"

type Auth struct {
	Model

	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) (bool, error) {
	var auth Auth

	if err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error; err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}
