package biz

import (
	"go-fiber-jwt/app/models"
	"go-fiber-jwt/app/repo"
)

type userBiz struct{}

type iUserBiz interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
}

var UserBiz iUserBiz

func init() {
	UserBiz = &userBiz{}
}

func (userBiz) GetUserByEmail(email string) (*models.User, error) {
	if len(email) <= 0 {
		return nil, nil
	}
	return repo.UserRepository.GetUserByEmail(email)
}

func (userBiz) GetUserById(id int) (*models.User, error) {
	if id <= 0 {
		return nil, nil
	}
	return repo.UserRepository.GetUserById(id)
}
