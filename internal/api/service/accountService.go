package service

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
)

func NewAccountService(r domain.AccountRepository) domain.AccountService {
	return &accountService{
		accountRepository: r,
	}
}

type accountService struct {
	accountRepository domain.AccountRepository
}

// Implementation of 'GetAccountByID'.
func (s *accountService) GetAccountByID(ctx context.Context, accountID, userID uint) (*domain.Account, error) {
	return s.accountRepository.GetAccountByID(ctx, accountID, userID)
}

// Implementation of 'GetAccountsOutputDTO'.
func (s *accountService) GetAccountsOutputDTO(ctx context.Context, filter *filter.Filter, userID uint) (*dto.ItemsOutputDTO, error) {
	return s.accountRepository.GetAccountsOutputDTO(ctx, filter, userID)
}

// Implementation of 'CreateAccount'.
func (s *accountService) CreateAccount(ctx context.Context, datas *dto.AccountInputDTO, userID uint) (*domain.Account, error) {
	return s.accountRepository.CreateAccount(ctx, datas, userID)
}

// Implementation of 'UpdateAccount'.
func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account, datas *dto.AccountInputDTO, userID uint) (*domain.Account, error) {
	return s.accountRepository.UpdateAccount(ctx, account, datas, userID)
}

// Implementation of 'DeleteAccount'.
func (s *accountService) DeleteAccount(ctx context.Context, account *domain.Account, userID uint) error {
	return s.accountRepository.DeleteAccount(ctx, account, userID)
}
