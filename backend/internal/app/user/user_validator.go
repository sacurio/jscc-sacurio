package user

import "errors"

type (
	// UserValidator determines the methods to be implemented by userValidator struct.
	UserValidator interface {
		Validate(user *User) error
	}

	userValidator struct{}
)

// NewUserValidator returns a new instance of UserValidator.
func NewUserValidator() UserValidator {
	return &userValidator{}
}

// Validate probes if a user info is correct before to be sent to repository.
func (uv *userValidator) Validate(user *User) error {
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
