package service

import (
	"errors"

	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
)

type UserOrchestration struct {
	userDomainService *domainService.UserDomainService
	userRepo          domainRepo.UserRepository
	roleRepo          domainRepo.RoleRepository
}

func NewUserOrchestration(
	userDomainService *domainService.UserDomainService,
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
) *UserOrchestration {
	return &UserOrchestration{
		userDomainService: userDomainService,
		userRepo:          userRepo,
		roleRepo:          roleRepo,
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

type BatchAssignRoleRequest struct {
	UserIDs []uint
	RoleID  uint
}

func (s *UserOrchestration) CreateUser(req *CreateUserRequest) (*entity.User, error) {
	if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
		return nil, errors.New("role not found")
	}

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

	if req.RoleID > 0 {
		if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
			return errors.New("role not found")
		}
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

func (s *UserOrchestration) AssignRole(userID uint, roleID uint) error {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return errors.New("user not found")
	}
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return errors.New("role not found")
	}
	return s.userRepo.Update(&entity.User{ID: userID, RoleID: roleID})
}

func (s *UserOrchestration) BatchAssignRole(req *BatchAssignRoleRequest) error {
	if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
		return errors.New("role not found")
	}
	if len(req.UserIDs) == 0 {
		return errors.New("user ids is required")
	}
	return s.userRepo.BatchUpdateRole(req.UserIDs, req.RoleID)
}

func (s *UserOrchestration) GetUsersByRole(roleID uint) ([]entity.User, error) {
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return nil, errors.New("role not found")
	}
	return s.userRepo.FindByRoleID(roleID)
}

type MenuTreeNode struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Path     string          `json:"path"`
	Sort     int             `json:"sort"`
	Children []*MenuTreeNode `json:"children"`
}

func buildMenuTree(permissions []entity.Permission, parentID uint) []*MenuTreeNode {
	var nodes []*MenuTreeNode
	for _, p := range permissions {
		if p.ParentID == parentID && p.Type == 1 {
			node := &MenuTreeNode{
				ID:       p.ID,
				Name:     p.Name,
				Code:     p.Code,
				Path:     p.Path,
				Sort:     p.Sort,
				Children: buildMenuTree(permissions, p.ID),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *UserOrchestration) GetCurrentUserPermissions(userID uint) ([]string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.RoleID == 0 {
		return nil, nil
	}
	permissions, err := s.roleRepo.GetPermissions(user.RoleID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, len(permissions))
	for i, p := range permissions {
		codes[i] = p.Code
	}
	return codes, nil
}

func (s *UserOrchestration) GetCurrentUserMenus(userID uint) ([]*MenuTreeNode, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.RoleID == 0 {
		return nil, nil
	}
	permissions, err := s.roleRepo.GetPermissions(user.RoleID)
	if err != nil {
		return nil, err
	}
	menuPerms := make([]entity.Permission, 0, len(permissions))
	for _, p := range permissions {
		if p.Type == 1 {
			menuPerms = append(menuPerms, p)
		}
	}
	return buildMenuTree(menuPerms, 0), nil
}

type UpdateProfileRequest struct {
	Nickname string
	Email    string
	Phone    string
}

func (s *UserOrchestration) UpdateProfile(userID uint, req *UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	return s.userRepo.Update(user)
}

func (s *UserOrchestration) UpdateAvatar(userID uint, avatarPath string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Avatar = avatarPath
	return s.userRepo.Update(user)
}

type UpdateStatusRequest struct {
	UserID uint
	Status int
}

func (s *UserOrchestration) UpdateStatus(userID uint, status int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Status = status
	return s.userRepo.Update(user)
}

func (s *UserOrchestration) ImportUsers(users []*CreateUserRequest) error {
	for _, req := range users {
		if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
			return errors.New("role not found for username: " + req.Username)
		}
		user := &entity.User{
			Username: req.Username,
			Nickname: req.Nickname,
			Email:    req.Email,
			Phone:    req.Phone,
			RoleID:   req.RoleID,
			Status:   1,
		}
		if err := s.userDomainService.Create(user, req.Password); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserOrchestration) ExportUsers() ([]entity.User, error) {
	users, _, err := s.userRepo.FindByPage(1, 10000, nil)
	if err != nil {
		return nil, err
	}
	return users, nil
}
