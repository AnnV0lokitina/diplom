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
	) (bool, error)
	AddUserSession(
		ctx context.Context,
		sessionID string,
		login string,
	) error
	GetUserBySessionID(ctx context.Context, activeSessionID string) (*entity.User, error)
	AddOrder(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber) error
	GetUserOrders(ctx context.Context, user *entity.User) ([]*entity.Order, error)
	GetUserBalance(ctx context.Context, user *entity.User) (*entity.UserBalance, error)
	UserOrderWithdraw(
		ctx context.Context,
		user *entity.User,
		order *entity.Order,
		sum entity.PointValue,
	) error
	GetUserOrderByNumber(
		ctx context.Context,
		user *entity.User,
		orderNumber entity.OrderNumber,
	) (*entity.Order, error)
	GetUserWithdrawals(ctx context.Context, user *entity.User) ([]*entity.Withdrawal, error)
	AddOrderInfo(
		ctx context.Context,
		orderNumber entity.OrderNumber,
		status entity.OrderStatus,
		accrual entity.PointValue,
	) error
	GetOrdersListForRequest(ctx context.Context) ([]entity.OrderNumber, error)
}

type Service struct {
	repo          Repo
	jobCheckOrder chan *entity.JobCheckOrder
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
	auth, err := s.repo.AuthUser(ctx, user.Login, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !auth {
		log.Info("user not found")
		return nil, labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
	}
	err = s.repo.AddUserSession(ctx, user.ActiveSessionID, user.Login)
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
	order, err := s.repo.GetUserOrderByNumber(ctx, user, orderNumber)
	if err != nil {
		return err
	}
	pointValue := entity.NewPointValue(sum)
	return s.repo.UserOrderWithdraw(ctx, user, order, pointValue)
}
