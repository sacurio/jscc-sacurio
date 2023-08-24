package repository

import (
	"errors"

	"github.com/sacurio/jb-challenge/internal/app/model"
	"gorm.io/gorm"
)

type (
	// User defines the methods to be implemented by the user service.
	User interface {
		Create(user *model.User) error
		FindByUsernameAndPassword(username, pwd string) (*model.User, error)
	}

	user struct {
		db *gorm.DB
	}
)

// Newuser returns a new instance of User Repository.
func NewUser(db *gorm.DB) User {
	return &user{
		db: db,
	}
}

// Create a new User through the use of Gorm.
func (ur *user) Create(user *model.User) error {
	result := ur.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByUsernameAndPassword determines if an user exists on data source, based on its username and password
// value and returns a *User instance.
func (ur *user) FindByUsernameAndPassword(username, pwd string) (*model.User, error) {
	var user model.User
	result := ur.db.Where("username = ? AND pwd = ?", username, pwd).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
