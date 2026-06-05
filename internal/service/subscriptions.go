package service

import (
	"Subscribe-service/internal/models"
	"github.com/google/uuid"
)

type SubsRepository interface {
	CreateSubs(subs *models.Subscription) error
	GetSubsByUserId(id uuid.UUID) (*models.Subscription, error)
	DeleteSubs(subs *models.Subscription) error
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

	return s.repo.CreateSubs(&subs)
}

func (s *SubsService) DeleteSubs(subs models.Subscription) error {
	//!праверки

	return s.repo.CreateSubs(&subs)
}

func (s *SubsService) GetSubs(Subs int) (models.Subscription, error) {
	//!праверки

	subs, err := s.repo.GetSubsByUserId(SubsId)
	if err != nil {
		return models.Subscription{}, err
	}

	return *subs, nil
}
