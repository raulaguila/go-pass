package service

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
)

func NewPhoneService(r domain.PhoneRepository) domain.PhoneService {
	return &phoneService{
		phoneRepository: r,
	}
}

type phoneService struct {
	phoneRepository domain.PhoneRepository
}

// Implementation of 'GetPhoneByID'.
func (s *phoneService) GetPhoneByID(ctx context.Context, phoneID uint) (*domain.Phone, error) {
	return s.phoneRepository.GetPhoneByID(ctx, phoneID)
}

// Implementation of 'GetPhonesOutputDTO'.
func (s *phoneService) GetPhonesOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	return s.phoneRepository.GetPhonesOutputDTO(ctx, filter)
}

// Implementation of 'CreatePhone'.
func (s *phoneService) CreatePhone(ctx context.Context, datas *dto.PhoneInputDTO) (*domain.Phone, error) {
	return s.phoneRepository.CreatePhone(ctx, datas)
}

// Implementation of 'UpdatePhone'.
func (s *phoneService) UpdatePhone(ctx context.Context, phone *domain.Phone, datas *dto.PhoneInputDTO) error {
	return s.phoneRepository.UpdatePhone(ctx, phone, datas)
}

// Implementation of 'DeletePhone'.
func (s *phoneService) DeletePhone(ctx context.Context, phone *domain.Phone) error {
	return s.phoneRepository.DeletePhone(ctx, phone)
}
