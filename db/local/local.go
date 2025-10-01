package local

import (
	"context"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"formation/db"
	"formation/model"
)

var _ db.UserStore = (*LocalUserStore)(nil)

// ou
// var _ db.UserStore = &LocalUserStore{}

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

func (s *LocalUserStore) GetUserByID(ctx context.Context, UUID string) (*model.User, error) {
	var out model.User
	if err := s.db.WithContext(ctx).First(&out, "uuid = ?", UUID).Error; err != nil {
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

func (s *LocalUserStore) UpdateUser(ctx context.Context, uuid string, u *model.User) error {
	return s.db.WithContext(ctx).Save(u).Error
}

func (s *LocalUserStore) DeleteUser(ctx context.Context, UUID string) error {
	return s.db.WithContext(ctx).Delete(&model.User{UUID: UUID}).Error
}

func (s *LocalUserStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
