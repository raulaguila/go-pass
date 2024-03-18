package domain

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
	"github.com/raulaguila/go-pass/pkg/validator"
)

const (
	AccountTableName      string = "account"
	AccountMailsTableName string = "account_mail_history"
	secretToken           string = "yZx4wwFE8CoJrAdW2GmOsw"
)

type (
	AccountMailHistory struct {
		Base
		AccountID uint     `json:"-" gorm:"column:account_id;type:bigint;not null;index;" validate:"required,min=1"`
		Account   *Account `json:"account,omitempty"`
		MailID    uint     `json:"-" gorm:"column:mail_id;type:bigint;not null;index;" validate:"required,min=1"`
		Mail      *Account `json:"mail,omitempty"`
	}

	Account struct {
		Base
		Username string   `json:"username" gorm:"column:username;type:varchar(100);index;not null;" validate:"required,min=2"`
		Password string   `json:"-" gorm:"column:password;type:varchar(100);index;not null;" validate:"required,min=2"`
		SiteID   uint     `json:"-" gorm:"column:site_id;type:bigint;not null;index;" validate:"required,min=1"`
		Site     *Site    `json:"site,omitempty"`
		PhoneID  uint     `json:"-" gorm:"column:phone_id;type:bigint;not null;index;" validate:"required,min=1"`
		Phone    *Phone   `json:"phone,omitempty"`
		MailID   uint     `json:"-" gorm:"column:mail_id;type:bigint;index;default:null;"`
		Mail     *Account `json:"mail,omitempty"`
		UserID   uint     `json:"-" gorm:"column:user_id;type:bigint;not null;index;" validate:"required,min=1"`
		User     *User    `json:"-"`
	}

	AccountRepository interface {
		GetAccountByID(context.Context, uint, uint) (*Account, error)
		GetAccountsOutputDTO(context.Context, *filter.Filter, uint) (*dto.ItemsOutputDTO, error)
		CreateAccount(context.Context, *dto.AccountInputDTO, uint) (*Account, error)
		UpdateAccount(context.Context, *Account, *dto.AccountInputDTO, uint) (*Account, error)
		DeleteAccount(context.Context, *Account, uint) error
	}

	AccountService interface {
		GetAccountByID(context.Context, uint, uint) (*Account, error)
		GetAccountsOutputDTO(context.Context, *filter.Filter, uint) (*dto.ItemsOutputDTO, error)
		CreateAccount(context.Context, *dto.AccountInputDTO, uint) (*Account, error)
		UpdateAccount(context.Context, *Account, *dto.AccountInputDTO, uint) (*Account, error)
		DeleteAccount(context.Context, *Account, uint) error
	}
)

func (Account) TableName() string {
	return AccountTableName
}

func (AccountMailHistory) TableName() string {
	return AccountMailsTableName
}

func (s *Account) EncodePass() {
	pass := fmt.Sprintf("%v%v%v", time.Now().Nanosecond(), secretToken, s.Password)
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(pass)))
	base64.StdEncoding.Encode(base64Text, []byte(pass))
	s.Password = string(base64Text)
}

func (s *Account) DecodePass() string {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(s.Password)))
	n, _ := base64.StdEncoding.Decode(base64Text, []byte(s.Password))
	return strings.Split(string(base64Text[:n]), secretToken)[1]
}

func (s *Account) Bind(p *dto.AccountInputDTO) error {
	if p.Username != nil {
		s.Username = *p.Username
	}

	if p.Password != nil {
		s.Password = *p.Password
		s.EncodePass()
	}

	if p.SiteID != nil {
		s.SiteID = *p.SiteID
	}

	if p.PhoneID != nil {
		s.PhoneID = *p.PhoneID
	}

	if p.MailID != nil {
		s.MailID = *p.MailID
	}

	return validator.StructValidator.Validate(s)
}

func (s Account) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"username": s.Username,
		"password": s.Password,
		"site_id":  s.SiteID,
		"phone_id": s.PhoneID,
		"mail_id":  s.MailID,
	}
}
