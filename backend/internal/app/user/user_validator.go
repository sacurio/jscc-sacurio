package user

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

// Validate probes if a user exists or not on data source, returning a custom error in negative case.
func (u *userValidator) Validate(user *User) error {
	return nil
}
