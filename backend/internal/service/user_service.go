package service

import (
	"context"
	"errors"
	"strings"
	"log"

	"github.com/cin/daily-spark/internal/model"
	"github.com/cin/daily-spark/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

var (
	ErrAccountExists  = errors.New("account already exists")
	ErrAccountNotExists = errors.New("account not exists")
	ErrPasswordNotMatch = errors.New("password not match")
	ErrUserNotFound = errors.New("user not found")
	ErrPermissionDenied = errors.New("permission denied")
)

func (s *UserService) UserSignUp(ctx context.Context, account, password string)(*model.User, error) {
	account = strings.TrimSpace(account)
	if account == "" || password == "" {
		log.Printf("[UserSignUp] empty account or password")
		return nil, errors.New("account or password empty")
	}

	_, err := s.repo.QueryUserByAccount(ctx, account)
	if err == nil {
		log.Printf("[UserSignUp] account exists: account=%s", account)
		return nil, ErrAccountExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[UserSignUp] bcrypt error: %v", err)
		return nil,err
	}

	user, err := s.repo.CreateUser(ctx, &model.User{
		Account:  account,
		Password: string(hashed),
	})
	if err != nil {
		log.Printf("[UserSignUp] create error: %v", err)
		return nil, err
	}
	log.Printf("[UserSignUp] created: id=%s account=%s", user.ID, user.Account)
	return user, nil
}

func (s *UserService) UserSignIn(ctx context.Context, account, password string) (*model.User, error) {
	account = strings.TrimSpace(account)
	if account == "" || password == "" {
		log.Printf("[UserSignIn] empty account or password")
		return nil, errors.New("account or password empty")
	}

	user, err := s.repo.QueryUserByAccount(ctx, account)
	if err != nil {
		log.Printf("[UserSignIn] account not exists: %s", account)
		return nil, ErrAccountNotExists
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("[UserSignIn] password not match: %s", account)
		return nil, ErrPasswordNotMatch
	}

	log.Printf("[UserSignIn] success: id=%s account=%s", user.ID, user.Account)
	return user, nil
}

func (s *UserService) QueryAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.QueryAllUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, nickname *string, password *string) (*model.User, error) {
	user, err := s.repo.QueryUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UpdateUser] user not found: id=%s", id)
			return nil, ErrUserNotFound
		}
		log.Printf("[UpdateUser] query error: %v", err)
		return nil, err
	}

	if nickname != nil {
		log.Printf("[UpdateUser] set nickname: id=%s", id)
		user.Nickname = nickname
	}

	if password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[UpdateUser] bcrypt error: %v", err)
			return nil, err
		}
		log.Printf("[UpdateUser] set password: id=%s", id)
		user.Password = string(hashed)
	}
	updated, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("[UpdateUser] save error: %v", err)
		return nil, err
	}
	log.Printf("[UpdateUser] success: id=%s", updated.ID)
	return updated, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id, account string) error {
	account = strings.TrimSpace(account)
	if account == "" {
		log.Printf("[DeleteUser] empty operator account")
		return errors.New("account empty")
	}

	operator, err := s.repo.QueryUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DeleteUser] operator not found: %s", account)
			return ErrUserNotFound
		}
		log.Printf("[DeleteUser] query operator error: %v", err)
		return err
	}

	if !operator.IsAdmin {
		log.Printf("[DeleteUser] permission denied: operator=%s", account)
		return ErrPermissionDenied
	}

	user, err := s.repo.QueryUserByAccount(ctx, account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[DeleteUser] user not found: %s", account)
			return ErrUserNotFound
		}
		log.Printf("[DeleteUser] query user error: %v", err)
		return err
	}

	if err := s.repo.DeleteUserByID(ctx, user.ID); err != nil {
		log.Printf("[DeleteUser] delete error: id=%s err=%v", user.ID, err)
		return err
	}
	log.Printf("[DeleteUser] success: id=%s by operator=%s", user.ID, operator.Account)
	return nil
}
