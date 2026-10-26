package service

import (
	"github.com/dannyexmakina/go-stencil-laravalian/app/model"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) List() ([]model.User, error) {
	var users []model.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *UserService) Create(name, email string) (*model.User, error) {
	u := &model.User{Name: name, Email: email, Status: model.UserStatusActive}
	err := s.db.Create(u).Error
	return u, err
}
