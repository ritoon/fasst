# TP — GORM avec interface + implémentations mock & locale (SQLite)

Objectifs :

- Définir un contrat `UserStore` (CRUD) dans `db/db.go`.
- Implémenter une version **mock** stateful (in-memory) dans `db/mock`.
- Implémenter une version **locale** avec **GORM + SQLite** dans `db/local`.
- Utiliser le modèle `User` centralisé dans `model/user.go`.

---

## 0) Pré-requis & initialisation

```bash
mkdir -p tp-echo/{model,db/{mock,local}}
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/google/uuid
```

---

## 1) Modèle de domaine

`model/user.go` :

```go
package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Soft delete possible si besoin :
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

---

## 2) Contrat (interface) du dépôt

`db/db.go` :

```go
package db

import (
	"context"
	"tp-echo/model"
)

type UserStore interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) error
	DeleteUser(ctx context.Context, id string) error
	Close() error
}
```

---

## 3) Implémentation **mock** (stateful, in-memory)

`db/mock/mock.go` :

```go
package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"tp-echo/db"
	"tp-echo/model"
)

var _ db.UserStore = (*MockUserStore)(nil)

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
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	for _, existing := range m.users {
		if existing.Email == u.Email {
			return errors.New("email already exists")
		}
	}
	cp := *u
	m.users[u.ID] = &cp
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *u
	return &cp, nil
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockUserStore) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*model.User, 0, len(m.users))
	for _, u := range m.users {
		cp := *u
		out = append(out, &cp)
	}
	if offset > len(out) {
		return []*model.User{}, nil
	}
	end := offset + limit
	if limit <= 0 || end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}

func (m *MockUserStore) UpdateUser(ctx context.Context, u *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[u.ID]; !ok {
		return errors.New("not found")
	}
	cp := *u
	m.users[u.ID] = &cp
	return nil
}

func (m *MockUserStore) DeleteUser(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[id]; !ok {
		return errors.New("not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserStore) Close() error { return nil }
```

---

## 4) Implémentation **locale** (GORM + SQLite)

`db/local/local.go` (squelette à compléter) :

```go
package local

import (
	"context"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tp-echo/db"
	"tp-echo/model"
)

var _ db.UserStore = (*LocalUserStore)(nil)

type LocalUserStore struct {
	db *gorm.DB
}

func New(path string) (*LocalUserStore, error) {
	gdb, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := gdb.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	return &LocalUserStore{db: gdb}, nil
}

func (s *LocalUserStore) CreateUser(ctx context.Context, u *model.User) error {
	return s.db.WithContext(ctx).Create(u).Error
}

func (s *LocalUserStore) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var out model.User
	if err := s.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &out, nil
}

func (s *LocalUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var out model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&out).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &out, nil
}

func (s *LocalUserStore) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	var list []*model.User
	q := s.db.WithContext(ctx).Model(&model.User{})
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if err := q.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *LocalUserStore) UpdateUser(ctx context.Context, u *model.User) error {
	return s.db.WithContext(ctx).Save(u).Error
}

func (s *LocalUserStore) DeleteUser(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&model.User{ID: id}).Error
}

func (s *LocalUserStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
```

---

## 5) Exemple d’usage

`main.go` :

```go
package main

import (
	"context"
	"fmt"
	"os"

	"tp-echo/db"
	"tp-echo/db/local"
	"tp-echo/db/mock"
	"tp-echo/model"
)

func buildStore() (db.UserStore, error) {
	if os.Getenv("USE_MOCK") == "1" {
		return mock.New(), nil
	}
	return local.New("app.db")
}

func main() {
	store, err := buildStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	ctx := context.Background()
	_ = store.CreateUser(ctx, &model.User{Email: "alice@example.com", Name: "Alice"})
	u, _ := store.GetUserByEmail(ctx, "alice@example.com")
	fmt.Println("Created:", u.ID, u.Email, u.Name)
}
```

---

## 6) Exercices proposés

1. Compléter ou réécrire les méthodes `db/local/local.go` en suivant [la documentation GORM](https://gorm.io/docs/).
2. Ajouter un champ `Role` à `User`, faire une migration et tester les CRUD.
3. Gérer les erreurs de doublon d’email avec `ErrRecordNotFound` et contraintes uniques.
4. Ajouter une fonction qui crée un user et met à jour son nom dans une **transaction**.
5. Implémenter un filtre `name LIKE` et la pagination dans `ListUsers`.
6. (Optionnel) Ajouter un `BeforeCreate` hook pour générer un `UUID` si non fourni.

---

## 7) Références utiles

- [Documentation GORM](https://gorm.io/docs/)
- [Driver SQLite](https://gorm.io/docs/connecting_to_the_database.html#SQLite)
- [Transactions GORM](https://gorm.io/docs/transactions.html)
- [Hooks](https://gorm.io/docs/hooks.html)
- [Scopes & requêtes avancées](https://gorm.io/docs/scopes.html)
