package service

import (
	"errors"
	"fmt"

	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
)

type RoleOrchestration struct {
	roleRepo domainRepo.RoleRepository
	permRepo domainRepo.PermissionRepository
}

func NewRoleOrchestration(roleRepo domainRepo.RoleRepository, permRepo domainRepo.PermissionRepository) *RoleOrchestration {
	return &RoleOrchestration{roleRepo: roleRepo, permRepo: permRepo}
}

type CreateRoleRequest struct {
	Name        string
	Code        string
	Description string
}

type UpdateRoleRequest struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Status      int
}

type AssignPermissionsRequest struct {
	RoleID        uint
	PermissionIDs []uint
}

func (s *RoleOrchestration) CreateRole(req *CreateRoleRequest) (*entity.Role, error) {
	existing, _ := s.roleRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("role code already exists")
	}

	role := &entity.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      int(entity.RoleStatusActive),
	}
	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleOrchestration) UpdateRole(req *UpdateRoleRequest) error {
	role, err := s.roleRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("role not found")
	}

	existing, _ := s.roleRepo.FindByCode(req.Code)
	if existing != nil && existing.ID != req.ID {
		return errors.New("role code already exists")
	}

	role.Name = req.Name
	role.Code = req.Code
	role.Description = req.Description
	role.Status = req.Status
	return s.roleRepo.Update(role)
}

func (s *RoleOrchestration) DeleteRole(id uint) error {
	if _, err := s.roleRepo.FindByID(id); err != nil {
		return errors.New("role not found")
	}
	return s.roleRepo.Delete(id)
}

func (s *RoleOrchestration) GetRole(id uint) (*entity.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleOrchestration) ListRoles(page, pageSize int) ([]entity.Role, int64, error) {
	return s.roleRepo.FindByPage(page, pageSize)
}

func (s *RoleOrchestration) ListAllRoles() ([]entity.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *RoleOrchestration) AssignPermissions(req *AssignPermissionsRequest) error {
	if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
		return errors.New("role not found")
	}
	for _, pid := range req.PermissionIDs {
		if _, err := s.permRepo.FindByID(pid); err != nil {
			return errors.New("permission not found: " + fmt.Sprintf("%d", pid))
		}
	}
	return s.roleRepo.AssignPermissions(req.RoleID, req.PermissionIDs)
}

func (s *RoleOrchestration) GetRolePermissions(roleID uint) ([]entity.Permission, error) {
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return nil, errors.New("role not found")
	}
	return s.roleRepo.GetPermissions(roleID)
}
