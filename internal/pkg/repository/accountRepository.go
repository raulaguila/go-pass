package repository

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/internal/pkg/postgre"
	"github.com/raulaguila/go-pass/pkg/filter"
	"gorm.io/gorm"
)

func NewAccountRepository(postgres *gorm.DB) domain.AccountRepository {
	return &accountRepository{
		postgres: postgres,
	}
}

type accountRepository struct {
	postgres *gorm.DB
}

func (s *accountRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *accountRepository) countAccounts(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.Account{}).Count(&count).Error
}

func (s *accountRepository) listAccounts(postgres *gorm.DB) (*[]domain.Account, error) {
	accounts := &[]domain.Account{}
	return accounts, postgres.Preload(postgre.Site).Preload(postgre.Mail).Preload(postgre.PhoneOperator).Find(accounts).Error
}

func (s *accountRepository) GetAccountsOutputDTO(ctx context.Context, filter *filter.Filter, userID uint) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countAccounts(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listAccounts(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *accountRepository) GetAccountByID(ctx context.Context, accountID, userID uint) (*domain.Account, error) {
	account := &domain.Account{}
	return account, s.postgres.Preload(postgre.Site).Preload(postgre.Mail).Preload(postgre.PhoneOperator).WithContext(ctx).First(account, accountID).Error
}

func (s *accountRepository) addAccountEmailHistory(ctx context.Context, accountID, mailID uint) error {
	history := &domain.AccountMailHistory{AccountID: accountID, MailID: mailID}
	return s.postgres.WithContext(ctx).Create(history).Error
}

func (s *accountRepository) CreateAccount(ctx context.Context, datas *dto.AccountInputDTO, userID uint) (*domain.Account, error) {
	account := &domain.Account{}
	if err := account.Bind(datas); err != nil {
		return nil, err
	}

	if err := s.postgres.WithContext(ctx).Create(account).Error; err != nil {
		return nil, err
	}

	if account.MailID != 0 {
		s.addAccountEmailHistory(ctx, account.Id, account.MailID)
	}

	return s.GetAccountByID(ctx, account.Id, userID)
}

func (s *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account, datas *dto.AccountInputDTO, userID uint) (*domain.Account, error) {
	mailID := account.MailID
	if err := account.Bind(datas); err != nil {
		return nil, err
	}

	if err := s.postgres.WithContext(ctx).Model(account).Updates(account.ToMap()).Error; err != nil {
		return nil, err
	}

	if account.MailID != mailID {
		s.addAccountEmailHistory(ctx, account.Id, account.MailID)
	}

	return s.GetAccountByID(ctx, account.Id, userID)
}

func (s *accountRepository) DeleteAccount(ctx context.Context, account *domain.Account, userID uint) error {
	return s.postgres.WithContext(ctx).Where("(id, user_id) IN ?", [][]interface{}{{account.Id, userID}}).Delete(&domain.Account{}).Error
}
