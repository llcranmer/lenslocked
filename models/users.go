package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrorNotFound = errors.New("models: resource not found")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(connectionString string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

func (us UserService) ByID(id uint) (*User, error) {
	var user User

	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrorNotFound
	default:
		return nil, err
	}
}

func (us *UserService) Create(u *User) error {
	return us.db.Create(u).Error
}

func (us *UserService) GetDB() *gorm.DB {
	return us.db
}

func (us UserService) Close() error {
	return us.db.Close()
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
