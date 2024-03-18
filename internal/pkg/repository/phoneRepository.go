package repository

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/internal/pkg/postgre"
	"github.com/raulaguila/go-pass/pkg/filter"
	"gorm.io/gorm"
)

func NewPhoneRepository(postgres *gorm.DB) domain.PhoneRepository {
	return &phoneRepository{
		postgres: postgres,
	}
}

type phoneRepository struct {
	postgres *gorm.DB
}

func (s *phoneRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *phoneRepository) countPhones(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.Phone{}).Count(&count).Error
}

func (s *phoneRepository) listPhones(postgres *gorm.DB) (*[]domain.Phone, error) {
	phones := &[]domain.Phone{}
	return phones, postgres.Preload(postgre.Operator).Preload(postgre.Accounts).Find(phones).Error
}

func (s *phoneRepository) GetPhonesOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countPhones(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listPhones(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *phoneRepository) GetPhoneByID(ctx context.Context, phoneID uint) (*domain.Phone, error) {
	phone := &domain.Phone{}
	return phone, s.postgres.WithContext(ctx).Preload(postgre.Operator).Preload(postgre.Accounts).First(phone, phoneID).Error
}

func (s *phoneRepository) CreatePhone(ctx context.Context, datas *dto.PhoneInputDTO) (*domain.Phone, error) {
	phone := &domain.Phone{}
	if err := phone.Bind(datas); err != nil {
		return nil, err
	}

	if err := s.postgres.WithContext(ctx).Create(phone).Error; err != nil {
		return nil, err
	}

	return s.GetPhoneByID(ctx, phone.Id)
}

func (s *phoneRepository) UpdatePhone(ctx context.Context, phone *domain.Phone, datas *dto.PhoneInputDTO) error {
	if err := phone.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(phone).Updates(phone.ToMap()).Error
}

func (s *phoneRepository) DeletePhone(ctx context.Context, phone *domain.Phone) error {
	return s.postgres.WithContext(ctx).Delete(phone).Error
}
