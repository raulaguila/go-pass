package service

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
)

func NewOperatorService(r domain.OperatorRepository) domain.OperatorService {
	return &operatorService{
		operatorRepository: r,
	}
}

type operatorService struct {
	operatorRepository domain.OperatorRepository
}

// Implementation of 'GetOperatorByID'.
func (s *operatorService) GetOperatorByID(ctx context.Context, operatorID uint) (*domain.Operator, error) {
	return s.operatorRepository.GetOperatorByID(ctx, operatorID)
}

// Implementation of 'GetOperatorsOutputDTO'.
func (s *operatorService) GetOperatorsOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	return s.operatorRepository.GetOperatorsOutputDTO(ctx, filter)
}

// Implementation of 'CreateOperator'.
func (s *operatorService) CreateOperator(ctx context.Context, datas *dto.OperatorInputDTO) (*domain.Operator, error) {
	return s.operatorRepository.CreateOperator(ctx, datas)
}

// Implementation of 'UpdateOperator'.
func (s *operatorService) UpdateOperator(ctx context.Context, operator *domain.Operator, datas *dto.OperatorInputDTO) error {
	return s.operatorRepository.UpdateOperator(ctx, operator, datas)
}

// Implementation of 'DeleteOperator'.
func (s *operatorService) DeleteOperator(ctx context.Context, operator *domain.Operator) error {
	return s.operatorRepository.DeleteOperator(ctx, operator)
}
