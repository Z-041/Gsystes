package service

import (
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
)

type UserOrchestration struct {
	userDomainService *domainService.UserDomainService
	userRepo          domainRepo.UserRepository
}

func NewUserOrchestration(
	userDomainService *domainService.UserDomainService,
	userRepo domainRepo.UserRepository,
) *UserOrchestration {
	return &UserOrchestration{
		userDomainService: userDomainService,
		userRepo:          userRepo,
	}
}

type CreateUserRequest struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	RoleID   uint
}

type UpdateUserRequest struct {
	ID       uint
	Nickname string
	Email    string
	Phone    string
	RoleID   uint
	Status   int
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	User  *entity.User
	Token string
}

func (s *UserOrchestration) CreateUser(req *CreateUserRequest) (*entity.User, error) {
	user := &entity.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	}

	if err := s.userDomainService.Create(user, req.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserOrchestration) UpdateUser(req *UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.RoleID = req.RoleID
	user.Status = req.Status

	return s.userRepo.Update(user)
}

func (s *UserOrchestration) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserOrchestration) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserOrchestration) ListUsers(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error) {
	return s.userRepo.FindByPage(page, pageSize, conditions)
}

func (s *UserOrchestration) Login(req *LoginRequest, tokenGenerator func(userID uint, username string, roleID uint) (string, error)) (*LoginResponse, error) {
	user, err := s.userDomainService.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := tokenGenerator(user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserOrchestration) ChangePassword(userID uint, oldPassword, newPassword string) error {
	return s.userDomainService.UpdatePassword(userID, oldPassword, newPassword)
}
