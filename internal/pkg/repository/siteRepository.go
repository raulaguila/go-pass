package repository

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
	"gorm.io/gorm"
)

func NewSiteRepository(postgres *gorm.DB) domain.SiteRepository {
	return &siteRepository{
		postgres: postgres,
	}
}

type siteRepository struct {
	postgres *gorm.DB
}

func (s *siteRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *siteRepository) countSites(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.Site{}).Count(&count).Error
}

func (s *siteRepository) listSites(postgres *gorm.DB) (*[]domain.Site, error) {
	sites := &[]domain.Site{}
	return sites, postgres.Find(sites).Error
}

func (s *siteRepository) GetSitesOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countSites(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listSites(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *siteRepository) GetSiteByID(ctx context.Context, siteID uint) (*domain.Site, error) {
	site := &domain.Site{}
	return site, s.postgres.WithContext(ctx).First(site, siteID).Error
}

func (s *siteRepository) CreateSite(ctx context.Context, datas *dto.SiteInputDTO) (*domain.Site, error) {
	site := &domain.Site{}
	if err := site.Bind(datas); err != nil {
		return nil, err
	}

	return site, s.postgres.WithContext(ctx).Create(site).Error
}

func (s *siteRepository) UpdateSite(ctx context.Context, site *domain.Site, datas *dto.SiteInputDTO) error {
	if err := site.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(site).Updates(site.ToMap()).Error
}

func (s *siteRepository) DeleteSite(ctx context.Context, site *domain.Site) error {
	return s.postgres.WithContext(ctx).Delete(site).Error
}
