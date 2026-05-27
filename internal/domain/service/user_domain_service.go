package service

import (
	"errors"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type UserDomainService struct {
	userRepo repository.UserRepository
}

func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
	return &UserDomainService{userRepo: userRepo}
}

func (s *UserDomainService) Create(user *entity.User, plainPassword string) error {
	existing, _ := s.userRepo.FindByUsername(user.Username)
	if existing != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Status = int(entity.UserStatusActive)

	return s.userRepo.Create(user)
}

func (s *UserDomainService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserDomainService) ValidateCredentials(username, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status == int(entity.UserStatusInactive) {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
