package repo

import (
	"go-fiber-jwt/app/models"
	"go-fiber-jwt/infra"

	"gorm.io/gorm"
)

type userRepo struct{}

type iUserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
}

var UserRepository iUserRepository

func init() {
	UserRepository = &userRepo{}
}

func (*userRepo) GetUserByEmail(email string) (*models.User, error) {
	if len(email) <= 0 {
		return nil, nil
	}
	var user models.User
	err := infra.DB.Where("email = ?", email).
		Find(&user).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, nil
}

func (*userRepo) GetUserById(id int) (*models.User, error) {
	if id <= 0 {
		return nil, nil
	}
	var user models.User
	err := infra.DB.Where("id = ?", id).
		Find(&user).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, nil
}
