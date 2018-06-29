package models

import (
	"../hash"
	"../rand"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorNotFound        = errors.New("models: resource not found")
	ErrorInvalidID       = errors.New("models: id must be > 0")
	ErrorInvalidPassword = errors.New("models: incorrect password provided")
)

const (
	userPepper    = "asdasdfaljfl;kj3;io4uklfjalkja#$#%@#sd4rhp2o83urowhrup8234u"
	hmacSecretKey = "asdasdasd32948723uhqkuryo782643198%$@!%^&#!t!^&%#!T!&^%&^^!@f!@&^!#%!&"
)

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	Close() error

	DestructiveReset() error
	AutoMigrate() error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserGorm(connectionString string) (*userGorm, error) {
	hmac := hash.NewHMAC(hmacSecretKey)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	return &userGorm{
		db:   db,
		HMAC: hmac,
	}, nil
}

var _ UserDB = &userGorm{}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

type userGorm struct {
	db *gorm.DB
	hash.HMAC
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

func (uv *userValidator) ByID(id uint) (*User, error) {
	if id < 0 {
		return nil, ErrorInvalidID
	}

	return uv.ByID(id)
}

func NewUserService(connectionString string) (UserService, error) {
	ug, err := NewUserGorm(connectionString)
	if err != nil {
		return nil, err
	}

	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

func (ug *userGorm) Create(u *User) error {
	pwBytes := []byte(u.Password + userPepper)
	hashBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.PasswordHash = string(hashBytes)
	u.Password = ""

	if u.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		u.Remember = token
	}
	u.RememberHash = ug.Hash(u.Remember)

	return ug.db.Create(u).Error
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

func (ug *userGorm) Update(u *User) error {
	if u.Remember != "" {
		u.RememberHash = ug.Hash(u.Remember)
	}

	return ug.db.Save(u).Error
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password+userPepper),
	)
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

func (ug *userGorm) Delete(id uint) error {
	if id < 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(user).Error
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := fist(db, &user)

	return &user, err
}

func (ug userGorm) ByID(id uint) (*User, error) {
	var user User

	db := ug.db.Where("id = ?", id)
	err := fist(db, &user)

	return &user, err
}

func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := ug.Hash(token)

	db := ug.db.Where("remember_hash = ?", hashedToken)
	err := fist(db, &user)
	if err != nil {
		return nil, err
	}

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

func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
