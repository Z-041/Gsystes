package service

import (
	"errors"
	"testing"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type mockUserRepo struct {
	users map[uint]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*entity.User)}
}

func (m *mockUserRepo) Create(user *entity.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(id uint) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepo) FindByUsername(username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByPage(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) FindByRoleID(roleID uint) ([]entity.User, error) {
	return nil, nil
}

func (m *mockUserRepo) BatchUpdateRole(userIDs []uint, roleID uint) error {
	return nil
}

func (m *mockUserRepo) BatchCreate(users []*entity.User) error {
	return nil
}

func TestCreateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{
		Username: "testuser",
	}
	err := svc.Create(user, "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be set")
	}
	if user.Status != int(entity.UserStatusActive) {
		t.Fatalf("expected status %d, got %d", entity.UserStatusActive, user.Status)
	}
	if !utils.CheckPassword("password123", user.Password) {
		t.Fatal("expected password to be hashed correctly")
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user1 := &entity.User{Username: "testuser"}
	if err := svc.Create(user1, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user2 := &entity.User{Username: "testuser"}
	err := svc.Create(user2, "password456")
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

func TestValidateCredentials_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, err := svc.ValidateCredentials("testuser", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "testuser" {
		t.Fatalf("expected username testuser, got %s", result.Username)
	}
}

func TestValidateCredentials_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := svc.ValidateCredentials("testuser", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "oldpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := svc.UpdatePassword(user.ID, "oldpassword", "newpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(user.ID)
	if !utils.CheckPassword("newpassword", updated.Password) {
		t.Fatal("expected password to be updated")
	}
}

func TestUpdatePassword_WrongOldPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "oldpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err := svc.UpdatePassword(user.ID, "wrongpassword", "newpassword")
	if err == nil {
		t.Fatal("expected error for wrong old password")
	}
}
