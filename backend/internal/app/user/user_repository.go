package user

import "gorm.io/gorm"

type (
	// UserRepository defines the methods to be implemented by the user service.
	UserRepository interface {
		Create(user *User) error
		FindByUsernameAndPassword(username, pwd string) (*User, error)
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

// Create a new User through the use of Gorm.
func (ur *userRepository) Create(user *User) error {
	result := ur.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByUsernameAndPassword determines if an user exists on data source, based on its username and password
// value and returns a *User instance.
func (ur *userRepository) FindByUsernameAndPassword(username, pwd string) (*User, error) {
	var user User
	result := ur.db.Where("username = ? AND pwd = ?", username, pwd).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}
