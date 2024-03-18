package repository

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
	"gorm.io/gorm"
)

func NewOperatorRepository(postgres *gorm.DB) domain.OperatorRepository {
	return &operatorRepository{
		postgres: postgres,
	}
}

type operatorRepository struct {
	postgres *gorm.DB
}

func (s *operatorRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *operatorRepository) countOperators(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.Operator{}).Count(&count).Error
}

func (s *operatorRepository) listOperators(postgres *gorm.DB) (*[]domain.Operator, error) {
	operators := &[]domain.Operator{}
	return operators, postgres.Find(operators).Error
}

func (s *operatorRepository) GetOperatorsOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countOperators(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listOperators(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *operatorRepository) GetOperatorByID(ctx context.Context, operatorID uint) (*domain.Operator, error) {
	operator := &domain.Operator{}
	return operator, s.postgres.WithContext(ctx).First(operator, operatorID).Error
}

func (s *operatorRepository) CreateOperator(ctx context.Context, datas *dto.OperatorInputDTO) (*domain.Operator, error) {
	operator := &domain.Operator{}
	if err := operator.Bind(datas); err != nil {
		return nil, err
	}

	return operator, s.postgres.WithContext(ctx).Create(operator).Error
}

func (s *operatorRepository) UpdateOperator(ctx context.Context, operator *domain.Operator, datas *dto.OperatorInputDTO) error {
	if err := operator.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(operator).Updates(operator.ToMap()).Error
}

func (s *operatorRepository) DeleteOperator(ctx context.Context, operator *domain.Operator) error {
	return s.postgres.WithContext(ctx).Delete(operator).Error
}
