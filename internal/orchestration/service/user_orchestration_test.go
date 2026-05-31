package service

import (
	"testing"

	"github.com/gsystes/backend/internal/domain/entity"
	domainService "github.com/gsystes/backend/internal/domain/service"
)

type mockTokenService struct{}

func (m *mockTokenService) GenerateToken(userID uint, username string, roleID uint) (string, error) {
	return "mock-token", nil
}

func TestUserOrchestration_CreateUser_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "Admin", Code: "admin"})
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, err := svc.CreateUser(&CreateUserRequest{
		Username: "newuser",
		Password: "pass123456",
		Nickname: "新用户",
		RoleID:   1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user ID to be set")
	}
	if user.Username != "newuser" {
		t.Fatalf("expected username newuser, got %s", user.Username)
	}
}

func TestUserOrchestration_CreateUser_RoleNotFound(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	_, err := svc.CreateUser(&CreateUserRequest{
		Username: "nobody",
		Password: "pass123",
		RoleID:   999,
	})
	if err == nil {
		t.Fatal("expected error for non-existent role")
	}
}

func TestUserOrchestration_AssignRole_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R1", Code: "r1"})
	roleRepo.Create(&entity.Role{Name: "R2", Code: "r2"})
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, _ := svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})

	if err := svc.AssignRole(user.ID, 2); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserOrchestration_UpdateProfile_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R", Code: "r"})
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, _ := svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})

	err := svc.UpdateProfile(user.ID, &UpdateProfileRequest{
		Nickname: "新昵称",
		Email:    "new@email.com",
		Phone:    "13800138000",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := userRepo.FindByID(user.ID)
	if got.Nickname != "新昵称" {
		t.Fatalf("expected nickname 新昵称, got %s", got.Nickname)
	}
}

func TestUserOrchestration_UpdateStatus_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R", Code: "r"})
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, _ := svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})

	if err := svc.UpdateStatus(user.ID, 2); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := userRepo.FindByID(user.ID)
	if got.Status != int(entity.UserStatusInactive) {
		t.Fatalf("expected status %d, got %d", entity.UserStatusInactive, got.Status)
	}
}

func TestUserOrchestration_UpdateAvatar_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R", Code: "r"})
	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, _ := svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})

	if err := svc.UpdateAvatar(user.ID, "/avatars/test.png"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, _ := userRepo.FindByID(user.ID)
	if got.Avatar != "/avatars/test.png" {
		t.Fatalf("expected avatar /avatars/test.png, got %s", got.Avatar)
	}
}

func TestUserOrchestration_GetCurrentUserMenus(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R", Code: "r"})
	roleRepo.AssignPermissions(1, []uint{1, 2})

	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	user, err := svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	menus, err := svc.GetCurrentUserMenus(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if menus == nil {
		t.Fatal("expected non-nil menus")
	}
}

func TestUserOrchestration_GetCurrentUserPermissions(t *testing.T) {
	userRepo := newMockUserRepo()
	roleRepo := newMockRoleRepo()
	roleRepo.Create(&entity.Role{Name: "R", Code: "r"})
	roleRepo.AssignPermissions(1, []uint{1, 2})

	userSvc := domainService.NewUserDomainService(userRepo)
	svc := NewUserOrchestration(userSvc, userRepo, roleRepo, &mockTokenService{})

	_, _ = svc.CreateUser(&CreateUserRequest{
		Username: "u1", Password: "pass123", RoleID: 1,
	})

	codes, err := svc.GetCurrentUserPermissions(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(codes) != 2 {
		t.Fatalf("expected 2 permission codes, got %d", len(codes))
	}
}
