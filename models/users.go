package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorNotFound        = errors.New("models: resource not found")
	ErrorInvalidID       = errors.New("models: id must be > 0")
	ErrorInvalidPassword = errors.New("models: incorrect password provided")
)

const userPepper = "asdasdfaljfl;kj3;io4uklfjalkjrhp2o83urowhrup8234u"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

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

func (us *UserService) Create(u *User) error {
	pwBytes := []byte(u.Password + userPepper)
	hashBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.PasswordHash = string(hashBytes)
	u.Password = ""
	return us.db.Create(u).Error
}

func (us *UserService) GetDB() *gorm.DB {
	return us.db
}

func (us *UserService) Close() error {
	return us.db.Close()
}

func (us *UserService) Update(u *User) error {
	return us.db.Save(u).Error
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := fist(db, &user)

	return &user, err
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+userPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrorInvalidPassword
		default:
			return nil, err
		}
	}

	return user, nil
}

func (us *UserService) Delete(id uint) error {
	if id < 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(user).Error
}

func (us UserService) ByID(id uint) (*User, error) {
	var user User

	db := us.db.Where("id = ?", id)
	err := fist(db, &user)

	return &user, err
}

func fist(db *gorm.DB, u *User) error {
	err := db.First(u).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrorNotFound
	}
	return err
}

func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
