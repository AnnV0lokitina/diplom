package service

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	"github.com/AnnV0lokitina/diplom/internal/mock"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_RegisterUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var session string
	login := "zo1q1NxnSrwV123"
	password := "A9XRCzNV1yS3izU9wD"
	ctx := context.TODO()
	passwordHash := entity.CreatePasswordHash(password)

	expectedUser := &entity.User{
		Login: login,
	}

	mockRepo := mock.NewMockRepo(ctl)

	errLoginExist := labelError.NewLabelError(labelError.TypeConflict, errors.New("login exists"))

	gomock.InOrder(
		mockRepo.EXPECT().CreateUser(ctx, gomock.AssignableToTypeOf(session), login, passwordHash).Return(nil),
		mockRepo.EXPECT().CreateUser(ctx, gomock.AssignableToTypeOf(session), login, passwordHash).Return(errLoginExist),
	)

	service := NewService(mockRepo)
	user, err := service.RegisterUser(ctx, login, password)
	assert.NoError(t, err)
	assert.Equal(t, user.Login, expectedUser.Login)

	user, err = service.RegisterUser(ctx, login, password)
	assert.Equal(t, err, errLoginExist)
	assert.Nil(t, user)
}

func TestService_LoginUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var login string
	var passwordHash string

	ctx := context.TODO()
	userID := 1

	user := &entity.User{
		ID:              userID,
		Login:           "login",
		ActiveSessionID: "sessionID",
	}

	mockRepo := mock.NewMockRepo(ctl)

	errUserNotFound := labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
	errUserUnauth := labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
	errDef := errors.New("error")

	gomock.InOrder(
		mockRepo.EXPECT().AuthUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(passwordHash)).Return(0, errUserNotFound),
		mockRepo.EXPECT().AuthUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(passwordHash)).Return(0, errDef),
		mockRepo.EXPECT().AuthUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(passwordHash)).Return(userID, nil),
		mockRepo.EXPECT().AuthUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(passwordHash)).Return(userID, nil),
	)

	gomock.InOrder(
		mockRepo.EXPECT().AddUserSession(ctx, gomock.AssignableToTypeOf(user)).Return(errDef),
		mockRepo.EXPECT().AddUserSession(ctx, gomock.AssignableToTypeOf(user)).Return(nil),
	)

	service := NewService(mockRepo)
	resultUser, err := service.LoginUser(ctx, user.Login, "password")
	assert.Nil(t, resultUser)
	assert.Equal(t, err.Error(), errUserUnauth.Error())

	resultUser, err = service.LoginUser(ctx, user.Login, "password")
	assert.Nil(t, resultUser)
	assert.Equal(t, err.Error(), errDef.Error())

	resultUser, err = service.LoginUser(ctx, user.Login, "password")
	assert.Nil(t, resultUser)
	assert.Equal(t, err.Error(), errDef.Error())

	resultUser, err = service.LoginUser(ctx, user.Login, "password")
	assert.Equal(t, resultUser.ID, userID)
	assert.Nil(t, err)
}

func TestService_authorizeUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var sessionID string

	ctx := context.TODO()
	user := &entity.User{
		ID:              1,
		Login:           "login",
		ActiveSessionID: "sessionID",
	}

	mockRepo := mock.NewMockRepo(ctl)

	errUserNotFound := labelError.NewLabelError(labelError.TypeNotFound, errors.New("no registered user"))
	errUserUnauth := labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("user unauthorized"))
	errDef := errors.New("error")

	gomock.InOrder(
		mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errUserNotFound),
		mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errDef),
		mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil),
	)

	service := NewService(mockRepo)
	resultUser, err := service.authorizeUser(ctx, "sessionID")
	assert.Nil(t, resultUser)
	assert.Equal(t, err.Error(), errUserUnauth.Error())

	resultUser, err = service.authorizeUser(ctx, "sessionID")
	assert.Nil(t, resultUser)
	assert.Equal(t, err.Error(), errDef.Error())

	resultUser, err = service.authorizeUser(ctx, "sessionID")
	assert.Equal(t, resultUser, user)
	assert.Nil(t, err)
}

func TestService_AddNewOrder(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var sessionID string
	var orderNumber entity.OrderNumber
	var user *entity.User

	correctNum := "74528626868518"
	incorrectNum := "77745286268685181"

	ctx := context.TODO()

	mockRepo := mock.NewMockRepo(ctl)

	errOrderNumIncorrect := labelError.NewLabelError(labelError.TypeInvalidData, errors.New("order number incorrect"))
	errNumCreated := labelError.NewLabelError(labelError.TypeCreated, errors.New("number created"))
	errNumExisted := labelError.NewLabelError(labelError.TypeConflict, errors.New("number exists"))
	errDef := errors.New("error")

	mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil).AnyTimes()

	gomock.InOrder(
		mockRepo.EXPECT().AddOrder(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errNumCreated),
		mockRepo.EXPECT().AddOrder(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errNumExisted),
		mockRepo.EXPECT().AddOrder(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errDef),
		mockRepo.EXPECT().AddOrder(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(nil),
	)

	service := NewService(mockRepo)

	err := service.AddNewOrder(ctx, "sessionID", incorrectNum)
	assert.Equal(t, err.Error(), errOrderNumIncorrect.Error())

	err = service.AddNewOrder(ctx, "sessionID", correctNum)
	assert.Equal(t, err.Error(), errNumCreated.Error())

	err = service.AddNewOrder(ctx, "sessionID", correctNum)
	assert.Equal(t, err.Error(), errNumExisted.Error())

	err = service.AddNewOrder(ctx, "sessionID", correctNum)
	assert.Equal(t, err.Error(), errDef.Error())

	err = service.AddNewOrder(ctx, "sessionID", correctNum)
	assert.Nil(t, err)
}

func TestService_GetOrderList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var sessionID string
	var user *entity.User
	var userOrders []*entity.Order

	userOrders = append(userOrders, &entity.Order{})

	ctx := context.TODO()

	mockRepo := mock.NewMockRepo(ctl)

	errNotFound := labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	errDef := errors.New("error")

	mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil).AnyTimes()

	gomock.InOrder(
		mockRepo.EXPECT().GetUserOrders(ctx, gomock.AssignableToTypeOf(user)).Return(nil, errNotFound),
		mockRepo.EXPECT().GetUserOrders(ctx, gomock.AssignableToTypeOf(user)).Return(nil, errDef),
		mockRepo.EXPECT().GetUserOrders(ctx, gomock.AssignableToTypeOf(user)).Return(userOrders, nil),
	)

	service := NewService(mockRepo)

	orders, err := service.GetOrderList(ctx, "sessionID")
	assert.Equal(t, err.Error(), errNotFound.Error())
	assert.Equal(t, len(orders), 0)

	orders, err = service.GetOrderList(ctx, "sessionID")
	assert.Equal(t, err.Error(), errDef.Error())
	assert.Equal(t, len(orders), 0)

	orders, err = service.GetOrderList(ctx, "sessionID")
	assert.Nil(t, err)
	assert.NotEqual(t, len(orders), 0)
}

func TestService_GetUserBalance(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var sessionID string
	var user *entity.User

	userBalance := &entity.UserBalance{}

	ctx := context.TODO()
	mockRepo := mock.NewMockRepo(ctl)
	errDef := errors.New("error")

	mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil).AnyTimes()

	gomock.InOrder(
		mockRepo.EXPECT().GetUserBalance(ctx, gomock.AssignableToTypeOf(user)).Return(nil, errDef),
		mockRepo.EXPECT().GetUserBalance(ctx, gomock.AssignableToTypeOf(user)).Return(userBalance, nil),
	)

	service := NewService(mockRepo)

	balance, err := service.GetUserBalance(ctx, "sessionID")
	assert.Equal(t, err.Error(), errDef.Error())
	assert.Nil(t, balance)

	balance, err = service.GetUserBalance(ctx, "sessionID")
	assert.Nil(t, err)
	assert.Equal(t, balance, userBalance)
}

func TestService_GetUserWithdrawals(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var sessionID string
	var user *entity.User
	var userWithdrawals []*entity.Withdrawal

	userWithdrawals = append(userWithdrawals, &entity.Withdrawal{})

	ctx := context.TODO()

	mockRepo := mock.NewMockRepo(ctl)

	errNotFound := labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	errDef := errors.New("error")

	mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil).AnyTimes()

	gomock.InOrder(
		mockRepo.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(user)).Return(nil, errNotFound),
		mockRepo.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(user)).Return(nil, errDef),
		mockRepo.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(user)).Return(userWithdrawals, nil),
	)

	service := NewService(mockRepo)

	withdrawals, err := service.GetUserWithdrawals(ctx, "sessionID")
	assert.Equal(t, err.Error(), errNotFound.Error())
	assert.Equal(t, len(withdrawals), 0)

	withdrawals, err = service.GetUserWithdrawals(ctx, "sessionID")
	assert.Equal(t, err.Error(), errDef.Error())
	assert.Equal(t, len(withdrawals), 0)

	withdrawals, err = service.GetUserWithdrawals(ctx, "sessionID")
	assert.Nil(t, err)
	assert.NotEqual(t, len(withdrawals), 0)
}

//func TestService_UserOrderWithdraw(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//
//	var sessionID string
//	var orderNumber entity.OrderNumber
//	var user *entity.User
//
//	correctNum := "74528626868518"
//	incorrectNum := "77745286268685181"
//	sum := 10.12
//
//	ctx := context.TODO()
//
//	mockRepo := mock.NewMockRepo(ctl)
//
//	errOrderNumIncorrect := labelError.NewLabelError(labelError.TypeInvalidData, errors.New("order number incorrect"))
//	errNumCreated := labelError.NewLabelError(labelError.TypeCreated, errors.New("number created"))
//	errNumExisted := labelError.NewLabelError(labelError.TypeConflict, errors.New("number exists"))
//	errDef := errors.New("error")
//
//	mockRepo.EXPECT().GetUserBySessionID(ctx, gomock.AssignableToTypeOf(sessionID)).Return(user, nil).AnyTimes()
//
//	gomock.InOrder(
//		mockRepo.EXPECT().UserOrderWithdraw(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errNumCreated),
//		mockRepo.EXPECT().UserOrderWithdraw(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errNumExisted),
//		mockRepo.EXPECT().UserOrderWithdraw(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(errDef),
//		mockRepo.EXPECT().UserOrderWithdraw(ctx, gomock.AssignableToTypeOf(user), gomock.AssignableToTypeOf(orderNumber)).Return(nil),
//	)
//
//	service := NewService(mockRepo)
//
//	err := service.UserOrderWithdraw(ctx, "sessionID", incorrectNum, sum)
//	assert.Equal(t, err.Error(), errOrderNumIncorrect.Error())
//
//	err = service.UserOrderWithdraw(ctx, "sessionID", correctNum, sum)
//	assert.Equal(t, err.Error(), errNumCreated.Error())
//
//	err = service.UserOrderWithdraw(ctx, "sessionID", correctNum, sum)
//	assert.Equal(t, err.Error(), errNumExisted.Error())
//
//	err = service.UserOrderWithdraw(ctx, "sessionID", correctNum, sum)
//	assert.Equal(t, err.Error(), errDef.Error())
//
//	err = service.UserOrderWithdraw(ctx, "sessionID", correctNum, sum)
//	assert.Nil(t, err)
//}
