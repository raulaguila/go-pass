package domain

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
	"github.com/raulaguila/go-pass/pkg/validator"
)

const SiteTableName string = "site"

type (
	Site struct {
		Base
		Name string `json:"name" gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
		URL  string `json:"url" gorm:"column:url;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	SiteRepository interface {
		GetSiteByID(context.Context, uint) (*Site, error)
		GetSitesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateSite(context.Context, *dto.SiteInputDTO) (*Site, error)
		UpdateSite(context.Context, *Site, *dto.SiteInputDTO) error
		DeleteSite(context.Context, *Site) error
	}

	SiteService interface {
		GetSiteByID(context.Context, uint) (*Site, error)
		GetSitesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateSite(context.Context, *dto.SiteInputDTO) (*Site, error)
		UpdateSite(context.Context, *Site, *dto.SiteInputDTO) error
		DeleteSite(context.Context, *Site) error
	}
)

func (Site) TableName() string {
	return SiteTableName
}

func (s *Site) Bind(p *dto.SiteInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	if p.URL != nil {
		s.URL = *p.URL
	}

	return validator.StructValidator.Validate(s)
}

func (s Site) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"name": s.Name,
		"url":  s.URL,
	}
}
