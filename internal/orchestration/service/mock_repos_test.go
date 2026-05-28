package service

import (
	"errors"
	"sync"

	"github.com/gsystes/backend/internal/domain/entity"
)

type mockUserRepo struct {
	mu     sync.RWMutex
	users  map[uint]*entity.User
	nextID uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*entity.User)}
}

func (m *mockUserRepo) Create(user *entity.User) error {
	m.nextID++
	user.ID = m.nextID
	m.mu.Lock()
	m.users[user.ID] = user
	m.mu.Unlock()
	return nil
}

func (m *mockUserRepo) Update(user *entity.User) error {
	m.mu.Lock()
	m.users[user.ID] = user
	m.mu.Unlock()
	return nil
}

func (m *mockUserRepo) Delete(id uint) error {
	m.mu.Lock()
	delete(m.users, id)
	m.mu.Unlock()
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByUsername(username string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByPage(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	all := make([]entity.User, 0, len(m.users))
	for _, u := range m.users {
		all = append(all, *u)
	}
	total := int64(len(all))
	start := (page - 1) * pageSize
	if start >= len(all) {
		return nil, total, nil
	}
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

func (m *mockUserRepo) FindByRoleID(roleID uint) ([]entity.User, error) {
	var result []entity.User
	for _, u := range m.users {
		if u.RoleID == roleID {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (m *mockUserRepo) BatchUpdateRole(userIDs []uint, roleID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, uid := range userIDs {
		if u, ok := m.users[uid]; ok {
			u.RoleID = roleID
		}
	}
	return nil
}

func (m *mockUserRepo) BatchCreate(users []*entity.User) error {
	for _, u := range users {
		m.Create(u)
	}
	return nil
}

type mockRoleRepo struct {
	mu       sync.RWMutex
	roles    map[uint]*entity.Role
	assigned map[uint][]entity.Permission
	nextID   uint
}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{
		roles:    make(map[uint]*entity.Role),
		assigned: make(map[uint][]entity.Permission),
	}
}

func (m *mockRoleRepo) Create(role *entity.Role) error {
	m.nextID++
	role.ID = m.nextID
	m.mu.Lock()
	m.roles[role.ID] = role
	m.mu.Unlock()
	return nil
}

func (m *mockRoleRepo) Update(role *entity.Role) error {
	m.mu.Lock()
	m.roles[role.ID] = role
	m.mu.Unlock()
	return nil
}

func (m *mockRoleRepo) Delete(id uint) error {
	m.mu.Lock()
	delete(m.roles, id)
	m.mu.Unlock()
	return nil
}

func (m *mockRoleRepo) FindByID(id uint) (*entity.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.roles[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return r, nil
}

func (m *mockRoleRepo) FindByCode(code string) (*entity.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.roles {
		if r.Code == code {
			return r, nil
		}
	}
	return nil, errors.New("role not found")
}

func (m *mockRoleRepo) FindAll() ([]entity.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]entity.Role, 0, len(m.roles))
	for _, r := range m.roles {
		result = append(result, *r)
	}
	return result, nil
}

func (m *mockRoleRepo) FindByPage(page, pageSize int) ([]entity.Role, int64, error) {
	all, _ := m.FindAll()
	total := int64(len(all))
	start := (page - 1) * pageSize
	if start >= len(all) {
		return nil, total, nil
	}
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

func (m *mockRoleRepo) AssignPermissions(roleID uint, permissionIDs []uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	perms := make([]entity.Permission, len(permissionIDs))
	for i, pid := range permissionIDs {
		perms[i] = entity.Permission{ID: pid, Code: "perm_code_" + string(rune('a'+pid)), Type: 1}
	}
	m.assigned[roleID] = perms
	return nil
}

func (m *mockRoleRepo) GetPermissions(roleID uint) ([]entity.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	perms, ok := m.assigned[roleID]
	if !ok {
		return []entity.Permission{}, nil
	}
	return perms, nil
}

type mockPermRepo struct {
	mu     sync.RWMutex
	perms  map[uint]*entity.Permission
	nextID uint
}

func newMockPermRepo() *mockPermRepo {
	return &mockPermRepo{perms: make(map[uint]*entity.Permission)}
}

func (m *mockPermRepo) Create(p *entity.Permission) error {
	m.nextID++
	p.ID = m.nextID
	m.mu.Lock()
	m.perms[p.ID] = p
	m.mu.Unlock()
	return nil
}

func (m *mockPermRepo) Update(p *entity.Permission) error {
	m.mu.Lock()
	m.perms[p.ID] = p
	m.mu.Unlock()
	return nil
}

func (m *mockPermRepo) Delete(id uint) error {
	m.mu.Lock()
	delete(m.perms, id)
	m.mu.Unlock()
	return nil
}

func (m *mockPermRepo) FindByID(id uint) (*entity.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.perms[id]
	if !ok {
		return nil, errors.New("permission not found")
	}
	return p, nil
}

func (m *mockPermRepo) FindByCode(code string) (*entity.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.perms {
		if p.Code == code {
			return p, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockPermRepo) FindAll() ([]entity.Permission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]entity.Permission, 0, len(m.perms))
	for _, p := range m.perms {
		result = append(result, *p)
	}
	return result, nil
}

func (m *mockPermRepo) FindByPage(page, pageSize int) ([]entity.Permission, int64, error) {
	return nil, 0, nil
}

func (m *mockPermRepo) FindByRoleID(roleID uint) ([]entity.Permission, error) {
	return nil, nil
}
