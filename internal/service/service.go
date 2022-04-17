package service

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
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
	AddOrder(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber) error
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

func (s *Service) AddNewOrder(ctx context.Context, sessionID string, orderNumber entity.OrderNumber) error {
	user, err := s.repo.GetUserBySessionID(ctx, sessionID)
	if err != nil {
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
			log.Info("user not found")
			return labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
		}
		return err
	}
	err = s.repo.AddOrder(ctx, user, orderNumber)
	if err != nil {
		return err
	}
	return nil
}
