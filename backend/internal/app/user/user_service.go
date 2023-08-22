package user

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
	return nil
}

// ValidateUser determines if an User exists on data source.
func (us *userService) ValidateUser(username, password string) (bool, error) {
	return false, nil
}
