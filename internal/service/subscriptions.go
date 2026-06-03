package service

import (
	"Subscribe-service/internal/models"
)

type TeamRepository interface {
	CreateSubs(team *models.Subscription) error
	GetSubsByUserId(name string) (*models.Subscription, error)
}

type SubsService struct {
	repo SubsRepository
}

func NewSubsService(repo SubsRepository) *SubsService {
	return &SubsService{
		repo: repo,
	}
}

func (s *SubsService) CreateSubs(subs models.Subscription) error {
	//!праверки

	return s.repo.CreateSubs(&team)
}

func (s *TeamService) GetSubs(Subs int) (models.Subscription, error) {
	//!праверки

	subs, err := s.repo.GetSubsByUserId(SubsId)
	if err != nil {
		return models.Subscription{}, err
	}

	return *subs, nil
}
