package domain

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
	"github.com/raulaguila/go-pass/pkg/validator"
)

const (
	PhoneTableName    string = "phone"
	OperatorTableName string = "operator"
)

type (
	Operator struct {
		Base
		Name string `json:"name" gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	Phone struct {
		Base
		Number     string    `json:"number" gorm:"column:number;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
		OperatorID uint      `json:"-" gorm:"column:operator_id;type:bigint;not null;index;" validate:"required,min=1"`
		Operator   *Operator `json:"operator,omitempty"`
	}

	PhoneRepository interface {
		GetPhoneByID(context.Context, uint) (*Phone, error)
		GetPhonesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreatePhone(context.Context, *dto.PhoneInputDTO) (*Phone, error)
		UpdatePhone(context.Context, *Phone, *dto.PhoneInputDTO) error
		DeletePhone(context.Context, *Phone) error
	}

	OperatorRepository interface {
		GetOperatorByID(context.Context, uint) (*Operator, error)
		GetOperatorsOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateOperator(context.Context, *dto.OperatorInputDTO) (*Operator, error)
		UpdateOperator(context.Context, *Operator, *dto.OperatorInputDTO) error
		DeleteOperator(context.Context, *Operator) error
	}

	PhoneService interface {
		GetPhoneByID(context.Context, uint) (*Phone, error)
		GetPhonesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreatePhone(context.Context, *dto.PhoneInputDTO) (*Phone, error)
		UpdatePhone(context.Context, *Phone, *dto.PhoneInputDTO) error
		DeletePhone(context.Context, *Phone) error
	}

	OperatorService interface {
		GetOperatorByID(context.Context, uint) (*Operator, error)
		GetOperatorsOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateOperator(context.Context, *dto.OperatorInputDTO) (*Operator, error)
		UpdateOperator(context.Context, *Operator, *dto.OperatorInputDTO) error
		DeleteOperator(context.Context, *Operator) error
	}
)

func (Phone) TableName() string {
	return PhoneTableName
}

func (Operator) TableName() string {
	return OperatorTableName
}

func (s *Phone) Bind(p *dto.PhoneInputDTO) error {
	if p.Number != nil {
		s.Number = *p.Number
	}

	if p.OperatorID != nil {
		s.OperatorID = *p.OperatorID
	}

	return validator.StructValidator.Validate(s)
}

func (s *Operator) Bind(p *dto.OperatorInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	return validator.StructValidator.Validate(s)
}

func (s Phone) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"number":      s.Number,
		"operator_id": s.OperatorID,
	}
}

func (s Operator) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"name": s.Name,
	}
}
