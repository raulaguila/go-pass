package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/internal/pkg/postgre"
	"github.com/raulaguila/go-pass/pkg/filter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserRepository(postgres *gorm.DB) domain.UserRepository {
	return &userRepository{
		postgres: postgres,
	}
}

type userRepository struct {
	postgres *gorm.DB
}

func (s *userRepository) applyFilter(ctx context.Context, filter *filter.UserFilter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	if filter.ProfileID != 0 {
		postgres = postgres.Where(domain.UserTableName+".profile_id = ?", filter.ProfileID)
	}
	postgres = postgres.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.UserTableName))
	postgres = filter.ApplySearchLike(postgres, domain.UserTableName+".name", domain.UserTableName+".mail", domain.ProfileTableName+".name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *userRepository) countUsers(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.User{}).Count(&count).Error
}

func (s *userRepository) listUsers(postgres *gorm.DB) (*[]domain.User, error) {
	users := &[]domain.User{}
	return users, postgres.Preload(postgre.ProfilePermission).Find(users).Error
}

func (s *userRepository) GetUsersOutputDTO(ctx context.Context, filter *filter.UserFilter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countUsers(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listUsers(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *userRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	user := &domain.User{}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).First(user, userID).Error
}

func (s *userRepository) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	user := &domain.User{Email: mail}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := &domain.User{Token: &token}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) CreateUser(ctx context.Context, datas *dto.UserInputDTO) (*domain.User, error) {
	user := &domain.User{New: true}
	if err := user.Bind(datas); err != nil {
		return nil, err
	}

	if err := s.postgres.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return s.GetUserByID(ctx, user.Id)
}

func (s *userRepository) UpdateUser(ctx context.Context, user *domain.User, datas *dto.UserInputDTO) error {
	if err := user.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(user).Updates(user.ToMap()).Error
}

func (s *userRepository) DeleteUser(ctx context.Context, user *domain.User) error {
	return s.postgres.WithContext(ctx).Delete(user).Error
}

func (s *userRepository) ResetUser(ctx context.Context, user *domain.User) error {
	user.Password = nil
	user.Token = nil
	user.New = true

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}

func (s *userRepository) PasswordUser(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	user.New = false
	user.Token = new(string)
	*user.Token = uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(*pass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = new(string)
	*user.Password = string(hash)

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}
