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

// 这里很关键，使用组合模式，并在后续的实现中，将数据库操作对象 *gorm.DB 作为参数传入实现解耦
type userRepository struct {
	db *gorm.DB
}

// 这里很关键，将 *gorm.DB 作为参数传入实现解耦
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
