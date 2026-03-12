package service

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
)

type CampaignService interface {
	CreateCampaign(req models.CreateCampaignRequest) (*models.Campaign, error)
	GetCampaign(id uint) (*models.Campaign, error)
}

type campaignService struct {
	repo repository.CampaignRepository
}

func NewCampaignService(r repository.CampaignRepository) CampaignService {
	return &campaignService{repo: r}
}

func (s *campaignService) CreateCampaign(req models.CreateCampaignRequest) (*models.Campaign, error) {
	campaign := &models.Campaign{
		Subject:     req.Subject,
		Body:        req.Body,
		Status:      models.CampaignStatusPending,
		TotalEmails: len(req.Recipients),
	}

	err := s.repo.CreateWithOutbox(campaign, req.Recipients)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *campaignService) GetCampaign(id uint) (*models.Campaign, error) {
	return s.repo.FindByID(id)
}
