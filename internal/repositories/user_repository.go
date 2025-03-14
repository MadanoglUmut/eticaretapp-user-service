package repositories

import (
	"UserService/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db}
}

func (r *UserRepository) GetUser(userId int) (models.Users, error) {

	var user models.Users

	if err := r.db.Table("users").Debug().Where("id = ?", userId).First(&user).Error; err != nil {

		return models.Users{}, err

	}

	return user, nil

}

func (r *UserRepository) GetUserByEmail(userEmail string) (models.Users, error) {

	var user models.Users

	if err := r.db.Table("users").Debug().Where("email = ?", userEmail).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(createdUser models.CreateUsers) (models.Users, error) {

	user := models.Users{
		Email:    createdUser.Email,
		Password: createdUser.Password,
		Isim:     createdUser.Isim,
		Soyisim:  createdUser.Soyisim,
	}

	if err := r.db.Table("users").Debug().Create(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil

}

func (r *UserRepository) UpdateUser(userId int, updatedUser models.UpdateUsers) (models.Users, error) {

	var user models.Users

	if err := r.db.Table("users").Debug().Where("id = ?", userId).First(&user).Error; err != nil {

		return models.Users{}, err

	}

	user.Password = updatedUser.Password
	user.Isim = updatedUser.Isim
	user.Soyisim = updatedUser.Soyisim

	if err := r.db.Table("users").Save(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil

}

func (r *UserRepository) DeleteUser(userId int) error {

	result := r.db.Table("users").Debug().Where("id = ?", userId).Delete(&models.Users{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil

}
