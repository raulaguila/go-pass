package service

import (
	"context"

	"github.com/raulaguila/go-pass/internal/pkg/domain"
	"github.com/raulaguila/go-pass/internal/pkg/dto"
	"github.com/raulaguila/go-pass/pkg/filter"
)

func NewSiteService(r domain.SiteRepository) domain.SiteService {
	return &siteService{
		siteRepository: r,
	}
}

type siteService struct {
	siteRepository domain.SiteRepository
}

// Implementation of 'GetSiteByID'.
func (s *siteService) GetSiteByID(ctx context.Context, siteID uint) (*domain.Site, error) {
	return s.siteRepository.GetSiteByID(ctx, siteID)
}

// Implementation of 'GetSitesOutputDTO'.
func (s *siteService) GetSitesOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	return s.siteRepository.GetSitesOutputDTO(ctx, filter)
}

// Implementation of 'CreateSite'.
func (s *siteService) CreateSite(ctx context.Context, datas *dto.SiteInputDTO) (*domain.Site, error) {
	return s.siteRepository.CreateSite(ctx, datas)
}

// Implementation of 'UpdateSite'.
func (s *siteService) UpdateSite(ctx context.Context, site *domain.Site, datas *dto.SiteInputDTO) error {
	return s.siteRepository.UpdateSite(ctx, site, datas)
}

// Implementation of 'DeleteSite'.
func (s *siteService) DeleteSite(ctx context.Context, site *domain.Site) error {
	return s.siteRepository.DeleteSite(ctx, site)
}
