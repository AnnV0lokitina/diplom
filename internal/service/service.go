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
	AuthUser(
		ctx context.Context,
		login string,
		passwordHash string,
	) (int, error)
	AddUserSession(ctx context.Context, user *entity.User) error
	GetUserBySessionID(ctx context.Context, activeSessionID string) (*entity.User, error)
	AddOrder(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber) error
	GetUserOrders(ctx context.Context, user *entity.User) ([]*entity.Order, error)
	GetUserBalance(ctx context.Context, user *entity.User) (*entity.UserBalance, error)
	UserOrderWithdraw(
		ctx context.Context,
		user *entity.User,
		orderNumber entity.OrderNumber,
		sum entity.PointValue,
	) error
	GetUserWithdrawals(ctx context.Context, user *entity.User) ([]*entity.Withdrawal, error)
	AddOrderInfo(ctx context.Context, orderInfo *entity.OrderUpdateInfo) error
	GetOrdersListForRequest(ctx context.Context) ([]entity.OrderNumber, error)
}

type AccrualSystem interface {
	GetOrderInfo(number entity.OrderNumber) (*entity.OrderUpdateInfo, error)
}

type JobCheckOrder struct {
	Number entity.OrderNumber
}

type Service struct {
	repo          Repo
	accrualSystem AccrualSystem
	jobCheckOrder chan *JobCheckOrder
}

func NewService(repo Repo, accrualSystem AccrualSystem) *Service {
	return &Service{
		repo:          repo,
		accrualSystem: accrualSystem,
	}
}

func (s *Service) RegisterUser(ctx context.Context, login string, password string) (*entity.User, error) {
	passwordHash := entity.CreatePasswordHash(password)
	sessionID, err := entity.GenerateSessionID()
	if err != nil {
		return nil, err
	}
	err = s.repo.CreateUser(ctx, sessionID, login, passwordHash)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		ActiveSessionID: sessionID,
		Login:           login,
	}
	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, login string, password string) (*entity.User, error) {
	passwordHash := entity.CreatePasswordHash(password)
	sessionID, err := entity.GenerateSessionID()
	if err != nil {
		return nil, err
	}
	userID, err := s.repo.AuthUser(ctx, login, passwordHash)
	if err != nil {
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
			log.Info("user not found")
			return nil, labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
		}
		return nil, err
	}
	user := &entity.User{
		ID:              userID,
		Login:           login,
		ActiveSessionID: sessionID,
	}
	err = s.repo.AddUserSession(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) authorizeUser(ctx context.Context, sessionID string) (*entity.User, error) {
	user, err := s.repo.GetUserBySessionID(ctx, sessionID)
	if err != nil {
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
			log.Info("user not found")
			return nil, labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) AddNewOrder(ctx context.Context, sessionID string, num string) error {
	user, err := s.authorizeUser(ctx, sessionID)
	if err != nil {
		return err
	}
	orderNumber := entity.OrderNumber(num)
	if !orderNumber.Valid() {
		log.Info("invalid order number")
		return labelError.NewLabelError(labelError.TypeInvalidData, errors.New("order number incorrect"))
	}
	return s.repo.AddOrder(ctx, user, orderNumber)
}

func (s *Service) GetOrderList(ctx context.Context, sessionID string) ([]*entity.Order, error) {
	user, err := s.authorizeUser(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserOrders(ctx, user)
}

func (s *Service) GetUserBalance(ctx context.Context, sessionID string) (*entity.UserBalance, error) {
	user, err := s.authorizeUser(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserBalance(ctx, user)
}

func (s *Service) GetUserWithdrawals(ctx context.Context, sessionID string) ([]*entity.Withdrawal, error) {
	user, err := s.authorizeUser(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserWithdrawals(ctx, user)
}

func (s *Service) UserOrderWithdraw(ctx context.Context, sessionID string, num string, sum float64) error {
	user, err := s.authorizeUser(ctx, sessionID)
	if err != nil {
		return err
	}
	orderNumber := entity.OrderNumber(num)
	if !orderNumber.Valid() {
		log.Info("invalid order number")
		return labelError.NewLabelError(labelError.TypeInvalidData, errors.New("order number incorrect"))
	}
	pointValue := entity.NewPointValue(sum)
	return s.repo.UserOrderWithdraw(ctx, user, orderNumber, pointValue)
}
