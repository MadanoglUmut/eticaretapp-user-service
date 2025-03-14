package repositories

import (
	"UserService/internal/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {

	fmt.Println("Hello", db)

	userRepository := NewUserRepository(db)

	createdUser := models.CreateUsers{
		Email:    "test@gmail.com",
		Password: "test123",
		Isim:     "Test",
		Soyisim:  "Test",
	}

	t.Run("CreateUser", func(t *testing.T) {

		user, err := userRepository.CreateUser(createdUser)

		assert.Nil(t, err)

		assert.Equal(t, createdUser.Email, user.Email)
	})

	t.Run("CreateUserError", func(t *testing.T) {

		createdUser := models.CreateUsers{
			Email:    "ahmet.yilmaz@example.com",
			Password: "test123",
			Isim:     "Test",
			Soyisim:  "Test",
		}

		_, err := userRepository.CreateUser(createdUser)

		assert.Error(t, err)

	})

	t.Run("GetUserByEmail", func(t *testing.T) {

		fetchedUser, err := userRepository.GetUserByEmail(createdUser.Email)

		assert.Nil(t, err)

		assert.Equal(t, createdUser.Email, fetchedUser.Email)

	})

	t.Run("GetUserByEmailError", func(t *testing.T) {

		fetchedUser, err := userRepository.GetUserByEmail(createdUser.Email)

		assert.Nil(t, err)

		assert.NotEqual(t, "ahmet.yilmaz@example.com", fetchedUser.Email)

	})

	t.Run("UpdateUser", func(t *testing.T) {

		fetchedUser, err := userRepository.GetUserByEmail(createdUser.Email)

		assert.Nil(t, err)

		updatedUser := models.UpdateUsers{Password: "123", Isim: "Guncel", Soyisim: "Guncel"}

		user, err := userRepository.UpdateUser(fetchedUser.ID, updatedUser)

		assert.Nil(t, err)

		assert.Equal(t, updatedUser.Password, user.Password)
		assert.Equal(t, updatedUser.Soyisim, user.Soyisim)

	})

	t.Run("DeleteUser", func(t *testing.T) {

		fetchedUser, err := userRepository.GetUserByEmail(createdUser.Email)

		assert.Nil(t, err)

		err = userRepository.DeleteUser(fetchedUser.ID)

		assert.Nil(t, err)

		_, err = userRepository.GetUser(fetchedUser.ID)

		assert.Error(t, err)

	})

	t.Run("DeleteUserError", func(t *testing.T) {

		err := userRepository.DeleteUser(999)

		assert.Equal(t, models.ErrRecordNotFound, err)

	})

}
