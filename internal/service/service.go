package service

import (
	"context"
	"github.com/AnnV0lokitina/diplom/internal/entity"
)

type Repo interface {
	Close(ctx context.Context) error
	CreateUser(
		ctx context.Context,
		sessionID string,
		login string,
		passwordHash string,
	) error
	UpdateUserSession(
		ctx context.Context,
		sessionID string,
		login string,
		passwordHash string,
	) error
	GetUserBySessionID(ctx context.Context, activeSessionID string) (*entity.User, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := entity.NewUser(login, password)
	if err != nil {
		return nil, err
	}
	err = s.repo.CreateUser(ctx, user.ActiveSessionID, user.Login, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, login string, password string) (*entity.User, error) {
	user, err := entity.NewUser(login, password)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpdateUserSession(ctx, user.ActiveSessionID, user.Login, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) AuthorizeUser(ctx context.Context, sessionID string) (*entity.User, error) {
	user, err := s.repo.GetUserBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
