package service

import (
	"testing"

	"github.com/gsystes/backend/internal/domain/entity"
)

func TestRoleOrchestration_CreateRole_Success(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	role, err := svc.CreateRole(&CreateRoleRequest{
		Name:        "测试角色",
		Code:        "test_role",
		Description: "测试用角色",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if role.ID == 0 {
		t.Fatal("expected role ID to be set")
	}
	if role.Status != int(entity.RoleStatusActive) {
		t.Fatalf("expected status active, got %d", role.Status)
	}
}

func TestRoleOrchestration_CreateRole_DuplicateCode(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	_, err := svc.CreateRole(&CreateRoleRequest{Name: "A", Code: "dup"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.CreateRole(&CreateRoleRequest{Name: "B", Code: "dup"})
	if err == nil {
		t.Fatal("expected error for duplicate code")
	}
}

func TestRoleOrchestration_UpdateRole_Success(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	role, _ := svc.CreateRole(&CreateRoleRequest{Name: "Old", Code: "old"})

	err := svc.UpdateRole(&UpdateRoleRequest{
		ID:          role.ID,
		Name:        "New",
		Code:        "new",
		Description: "updated",
		Status:      2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := rRepo.FindByID(role.ID)
	if updated.Name != "New" {
		t.Fatalf("expected name New, got %s", updated.Name)
	}
	if updated.Status != 2 {
		t.Fatalf("expected status 2, got %d", updated.Status)
	}
}

func TestRoleOrchestration_DeleteRole_Success(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	role, _ := svc.CreateRole(&CreateRoleRequest{Name: "X", Code: "x"})

	if err := svc.DeleteRole(role.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := rRepo.FindByID(role.ID)
	if err == nil {
		t.Fatal("expected role to be deleted")
	}
}

func TestRoleOrchestration_ListRoles(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	svc.CreateRole(&CreateRoleRequest{Name: "A", Code: "a"})
	svc.CreateRole(&CreateRoleRequest{Name: "B", Code: "b"})

	roles, total, err := svc.ListRoles(1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
}

func TestRoleOrchestration_AssignPermissions_Success(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	pRepo.Create(&entity.Permission{Name: "P1", Code: "p1"})
	pRepo.Create(&entity.Permission{Name: "P2", Code: "p2"})
	role, _ := svc.CreateRole(&CreateRoleRequest{Name: "R", Code: "r"})

	err := svc.AssignPermissions(&AssignPermissionsRequest{
		RoleID:        role.ID,
		PermissionIDs: []uint{1, 2},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRoleOrchestration_AssignPermissions_PermissionNotFound(t *testing.T) {
	rRepo := newMockRoleRepo()
	pRepo := newMockPermRepo()
	svc := NewRoleOrchestration(rRepo, pRepo)

	role, _ := svc.CreateRole(&CreateRoleRequest{Name: "R", Code: "r"})

	err := svc.AssignPermissions(&AssignPermissionsRequest{
		RoleID:        role.ID,
		PermissionIDs: []uint{999},
	})
	if err == nil {
		t.Fatal("expected error for non-existent permission")
	}
}
