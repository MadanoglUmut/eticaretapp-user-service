package services

import "UserService/internal/models"

type userRepository interface {
	GetUser(userId int) (models.Users, error)
	GetUserByEmail(userEmail string) (models.Users, error)
	CreateUser(createdUser models.CreateUsers) (models.Users, error)
	UpdateUser(userId int, updatedUser models.UpdateUsers) (models.Users, error)
	DeleteUser(userId int) error
}

type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUser(userId int) (models.Users, error) {
	return s.userRepository.GetUser(userId)
}

func (s *UserService) GetUserByEmail(userEmail string) (models.Users, error) {
	return s.userRepository.GetUserByEmail(userEmail)
}

func (s *UserService) CreateUser(createdUser models.CreateUsers) (models.Users, error) {
	return s.userRepository.CreateUser(createdUser)
}

func (s *UserService) UpdateUser(userId int, updatedUser models.UpdateUsers) (models.Users, error) {
	return s.userRepository.UpdateUser(userId, updatedUser)
}

func (s *UserService) Delete(userId int) error {
	return s.userRepository.DeleteUser(userId)
}
