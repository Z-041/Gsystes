package service

import (
	"errors"

	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
)

type PermissionOrchestration struct {
	permRepo domainRepo.PermissionRepository
}

func NewPermissionOrchestration(permRepo domainRepo.PermissionRepository) *PermissionOrchestration {
	return &PermissionOrchestration{permRepo: permRepo}
}

type CreatePermissionRequest struct {
	Name     string
	Code     string
	Type     int
	ParentID uint
	Path     string
	Method   string
	Sort     int
}

type UpdatePermissionRequest struct {
	ID       uint
	Name     string
	Code     string
	Type     int
	ParentID uint
	Path     string
	Method   string
	Sort     int
}

func (s *PermissionOrchestration) CreatePermission(req *CreatePermissionRequest) (*entity.Permission, error) {
	existing, _ := s.permRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("permission code already exists")
	}

	p := &entity.Permission{
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		ParentID: req.ParentID,
		Path:     req.Path,
		Method:   req.Method,
		Sort:     req.Sort,
	}
	if err := s.permRepo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PermissionOrchestration) UpdatePermission(req *UpdatePermissionRequest) error {
	p, err := s.permRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("permission not found")
	}

	existing, _ := s.permRepo.FindByCode(req.Code)
	if existing != nil && existing.ID != req.ID {
		return errors.New("permission code already exists")
	}

	p.Name = req.Name
	p.Code = req.Code
	p.Type = req.Type
	p.ParentID = req.ParentID
	p.Path = req.Path
	p.Method = req.Method
	p.Sort = req.Sort
	return s.permRepo.Update(p)
}

func (s *PermissionOrchestration) DeletePermission(id uint) error {
	if _, err := s.permRepo.FindByID(id); err != nil {
		return errors.New("permission not found")
	}
	return s.permRepo.Delete(id)
}

func (s *PermissionOrchestration) GetPermission(id uint) (*entity.Permission, error) {
	return s.permRepo.FindByID(id)
}

func (s *PermissionOrchestration) ListPermissions(page, pageSize int) ([]entity.Permission, int64, error) {
	return s.permRepo.FindByPage(page, pageSize)
}

func (s *PermissionOrchestration) ListAllPermissions() ([]entity.Permission, error) {
	return s.permRepo.FindAll()
}
