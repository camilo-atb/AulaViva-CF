package services

import (
	"context"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/repositories"
)

type AuthorizationService struct {
	authorizations repositories.AuthorizationRepository
}

func NewAuthorizationService(authorizations repositories.AuthorizationRepository) *AuthorizationService {
	return &AuthorizationService{authorizations: authorizations}
}

func (s *AuthorizationService) Get(ctx context.Context, token string) (dto.RequestResponse, error) {
	authorization, event, student, err := s.authorizations.GetByToken(ctx, token)
	if err != nil {
		return dto.RequestResponse{}, err
	}
	return dto.RequestResponse{Authorization: authorization, Event: event, Student: student}, nil
}

func (s *AuthorizationService) Respond(ctx context.Context, token string, input dto.SubmitResponseInput) (dto.RequestResponse, error) {
	if err := (entities.Authorization{}).ValidateResponse(input.Status); err != nil {
		return dto.RequestResponse{}, err
	}
	now := time.Now().UTC()
	if err := s.authorizations.Respond(ctx, token, input.Status, now); err != nil {
		return dto.RequestResponse{}, err
	}
	return s.Get(ctx, token)
}
