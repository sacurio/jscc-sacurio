package user

import (
	"crypto/md5"
	"fmt"
)

type (
	// UserService
	UserService interface {
		RegisterUser(username, email, password string) error
		ValidateUser(username, password string) (bool, error)
	}

	userService struct {
		userRepo      UserRepository
		userValidator UserValidator
	}
)

// NewService returns a new instance of UserService.
func NewUserService(userRepo UserRepository, userValidator UserValidator) UserService {
	return &userService{
		userRepo:      userRepo,
		userValidator: userValidator,
	}
}

// RegisterUser persists a new User to data source.
func (us *userService) RegisterUser(username, email, password string) error {
	// Validar el usuario
	newUser := &User{
		Username: username,
		Email:    email,
		Pwd:      password,
	}

	if err := us.userValidator.Validate(newUser); err != nil {
		return err
	}

	hashedPassword := md5.Sum([]byte(password))
	newUser.Pwd = fmt.Sprintf("%x", hashedPassword)

	if err := us.userRepo.Create(newUser); err != nil {
		return err
	}

	return nil
}

// ValidateUser determines if an User exists on data source.
func (us *userService) ValidateUser(username, password string) (bool, error) {
	hashedPassword := md5.Sum([]byte(password))
	strHashedPassword := fmt.Sprintf("%x", hashedPassword)

	user, err := us.userRepo.FindByUsernameAndPassword(username, strHashedPassword)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil // Usuario no encontrado
	}

	return true, nil
}
