package service

import (
	"Subscribe-service/internal/models"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("subscription not found")
	ErrInvalid  = errors.New("invalid subscription data")
)

type SubsRepository interface {
	CreateSubs(subs *models.Subscription) error
	GetSubs(userID uuid.UUID, serviceName string) (*models.Subscription, error)
	GetAllSubs() ([]models.Subscription, error)
	UpdateSubs(subs *models.Subscription) error
	DeleteSubs(userID uuid.UUID, serviceName string) error
	GetSubsSum(req models.AggregateRequest) (int, error)
}

type SubsService struct {
	repo SubsRepository
}

func NewSubsService(repo SubsRepository) *SubsService {
	return &SubsService{repo: repo}
}

func (s *SubsService) CreateSubs(subs models.Subscription) error {
	if err := validateSubs(subs); err != nil {
		return err
	}

	return s.repo.CreateSubs(&subs)
}

func (s *SubsService) GetSubs(userID uuid.UUID, serviceName string) (models.Subscription, error) {
	sub, err := s.repo.GetSubs(userID, serviceName)
	if err != nil {
		return models.Subscription{}, err
	}

	if sub == nil {
		return models.Subscription{}, ErrNotFound
	}

	return *sub, nil
}

func (s *SubsService) UpdateSubs(subs models.Subscription) error {
	if err := validateSubs(subs); err != nil {
		return err
	}

	return s.repo.UpdateSubs(&subs)
}

func (s *SubsService) DeleteSubs(userID uuid.UUID, serviceName string) error {
	if serviceName == "" {
		return ErrInvalid
	}

	return s.repo.DeleteSubs(userID, serviceName)
}

func (s *SubsService) GetAllSubs() ([]models.Subscription, error) {
	return s.repo.GetAllSubs()
}

func (s *SubsService) GetSubsSum(req models.AggregateRequest) (int, error) {
	return s.repo.GetSubsSum(req)
}

func validateSubs(sub models.Subscription) error {
	if sub.UserID == uuid.Nil {
		return ErrInvalid
	}

	if strings.TrimSpace(sub.ServiceName) == "" {
		return ErrInvalid
	}

	if sub.Price <= 0 {
		return ErrInvalid
	}

	if strings.TrimSpace(sub.StartDate) == "" {
		return ErrInvalid
	}

	return nil
}
