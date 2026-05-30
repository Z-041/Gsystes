package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	"golang.org/x/sync/errgroup"
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
		return nil, fmt.Errorf("role not found: %w", err)
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
			return fmt.Errorf("role not found: %w", err)
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
		return fmt.Errorf("user not found: %w", err)
	}
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	return s.userRepo.Update(&entity.User{ID: userID, RoleID: roleID})
}

func (s *UserOrchestration) BatchAssignRole(req *BatchAssignRoleRequest) error {
	if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	if len(req.UserIDs) == 0 {
		return errors.New("user ids is required")
	}
	return s.userRepo.BatchUpdateRole(req.UserIDs, req.RoleID)
}

func (s *UserOrchestration) GetUsersByRole(roleID uint) ([]entity.User, error) {
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
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
	childrenMap := make(map[uint][]entity.Permission)
	roots := make([]entity.Permission, 0)

	for _, p := range permissions {
		if p.Type != 1 {
			continue
		}
		if p.ParentID == parentID {
			roots = append(roots, p)
		}
		childrenMap[p.ParentID] = append(childrenMap[p.ParentID], p)
	}

	var build func(pid uint) []*MenuTreeNode
	build = func(pid uint) []*MenuTreeNode {
		children, ok := childrenMap[pid]
		if !ok {
			return nil
		}
		nodes := make([]*MenuTreeNode, 0, len(children))
		for _, p := range children {
			node := &MenuTreeNode{
				ID:       p.ID,
				Name:     p.Name,
				Code:     p.Code,
				Path:     p.Path,
				Sort:     p.Sort,
				Children: build(p.ID),
			}
			nodes = append(nodes, node)
		}
		return nodes
	}

	return build(parentID)
}

func (s *UserOrchestration) GetCurrentUserPermissions(userID uint) ([]string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
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
		return nil, fmt.Errorf("user not found: %w", err)
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
	type validatedUser struct {
		entity *entity.User
		err    error
	}

	reqCh := make(chan *CreateUserRequest, len(users))
	resultCh := make(chan *validatedUser, len(users))
	g, ctx := errgroup.WithContext(context.Background())

	workerCount := 8
	for i := 0; i < workerCount; i++ {
		g.Go(func() error {
			for req := range reqCh {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
					resultCh <- &validatedUser{err: fmt.Errorf("role not found for %s: %w", req.Username, err)}
					continue
				}
				existing, _ := s.userRepo.FindByUsername(req.Username)
				if existing != nil {
					resultCh <- &validatedUser{err: fmt.Errorf("username already exists: %s", req.Username)}
					continue
				}
				hashedPassword, err := utils.HashPassword(req.Password)
				if err != nil {
					resultCh <- &validatedUser{err: err}
					continue
				}
				resultCh <- &validatedUser{
					entity: &entity.User{
						Username: req.Username,
						Password: hashedPassword,
						Nickname: req.Nickname,
						Email:    req.Email,
						Phone:    req.Phone,
						RoleID:   req.RoleID,
						Status:   1,
					},
				}
			}
			return nil
		})
	}

	go func() {
		for _, req := range users {
			reqCh <- req
		}
		close(reqCh)
	}()

	go func() {
		g.Wait()
		close(resultCh)
	}()

	var entities []*entity.User
	for r := range resultCh {
		if r.err != nil {
			return r.err
		}
		entities = append(entities, r.entity)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if len(entities) == 0 {
		return nil
	}

	return s.userRepo.BatchCreate(entities)
}
