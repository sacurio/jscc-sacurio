package user

import "gorm.io/gorm"

type (
	// UserRepository defines the methods to be implemented by the user service.
	UserRepository interface {
		Create(user *User) error
		FindByUsername(username string) (*User, error)
	}

	userRepository struct {
		db *gorm.DB
	}
)

// NewUserRepository returns a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// TODO: Probably be deleted
// Create...
func (r *userRepository) Create(user *User) error {
	return nil
}

// FindByUsername determines if an user exists on data source, based on its username value and returns a *User instance.
func (r *userRepository) FindByUsername(username string) (*User, error) {
	return nil, nil
}
