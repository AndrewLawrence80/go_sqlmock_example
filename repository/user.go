package repository

import (
	"example/entity"

	"gorm.io/gorm"
)

//go:generate gotests -exported -w user.go

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint64) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.User{}, id).Error
}
