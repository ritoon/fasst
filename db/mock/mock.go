package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"formation/db"
	"formation/model"
)

var _ db.UserStore = (*MockUserStore)(nil)

// ou
// var _ db.UserStore = &MockUserStore{}

type MockUserStore struct {
	mu    sync.RWMutex
	users map[string]*model.User // ID -> User
}

func New() *MockUserStore {
	return &MockUserStore{users: make(map[string]*model.User)}
}

func (m *MockUserStore) CreateUser(ctx context.Context, u *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.users {
		if existing.Email == u.Email {
			return errors.New("email already exists")
		}
	}
	u.UUID = uuid.NewString()
	m.users[u.UUID] = u
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, uuid string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[uuid]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockUserStore) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*model.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, u)
	}

	return out, nil
}

func (m *MockUserStore) UpdateUser(ctx context.Context, uuid string, u *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[uuid]; !ok {
		return errors.New("not found")
	}
	m.users[uuid] = u
	return nil
}

func (m *MockUserStore) DeleteUser(ctx context.Context, uuid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[uuid]; !ok {
		return errors.New("not found")
	}
	delete(m.users, uuid)
	return nil
}

func (m *MockUserStore) Close() error { return nil }
