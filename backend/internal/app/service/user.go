package service

import (
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/sacurio/jb-challenge/internal/app/model"
	repository "github.com/sacurio/jb-challenge/internal/app/repository/user"
)

type (
	// User
	User interface {
		Register(username, email, password string) error
		Validate(username, password string) (bool, error)
	}

	user struct {
		repository repository.User
	}
)

// NewService returns a new instance of UserService.
func NewService(repository repository.User) User {
	return &user{
		repository: repository,
	}
}

// Register persists a new User to data source.
func (us *user) Register(username, email, password string) error {
	newUser := &model.User{
		Username: username,
		Email:    email,
		Pwd:      password,
	}

	if err := validate(newUser); err != nil {
		return err
	}

	hashedPassword := md5.Sum([]byte(password))
	newUser.Pwd = fmt.Sprintf("%x", hashedPassword)

	if err := us.repository.Create(newUser); err != nil {
		return err
	}

	return nil
}

// ValidateUser determines if an User exists on data source.
func (us *user) Validate(username, password string) (bool, error) {
	hashedPassword := md5.Sum([]byte(password))
	strHashedPassword := fmt.Sprintf("%x", hashedPassword)

	user, err := us.repository.FindByUsernameAndPassword(username, strHashedPassword)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return true, nil
}

func validate(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Pwd == "" {
		return errors.New("password is required")
	}

	return nil
}
