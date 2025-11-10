package repository

import (
	"context"
	"github.com/cin/daily-spark/internal/model"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("[Repo.CreateUser] db error: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) QueryUserByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("account = ? and deleted_at is null", account).First(&user).Error; err != nil {
		log.Printf("[Repo.QueryUserByAccount] db error: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) QueryUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&user).Error; err != nil {
		log.Printf("[Repo.QueryUserByID] db error: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) QueryAllUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Where("deleted_at is null").Find(&users).Error; err != nil {
		log.Printf("[Repo.QueryAllUsers] db error: %v", err)
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Printf("[Repo.UpdateUser] db error: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&model.User{}).Error
	if err != nil {
		log.Printf("[Repo.DeleteUserByID] db error: id=%s err=%v", id, err)
	}
	return err
}
