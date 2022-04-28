package handler

import (
	"context"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"time"

	"github.com/AnnV0lokitina/diplom/internal/entity"
	"github.com/golang/mock/gomock"

	mock "github.com/AnnV0lokitina/diplom/internal/handler_mock"
)

func createMockService(ctl *gomock.Controller) Service {
	mockService := mock.NewMockService(ctl)

	ctx := context.TODO()

	var login string
	var password string
	var sessionID string
	var num string
	var sum float64

	user := &entity.User{
		Login:           "login",
		ActiveSessionID: "session",
	}

	errDef := errors.New("error")
	errConflict := labelError.NewLabelError(labelError.TypeConflict, errors.New("conflict"))
	errUnauth := labelError.NewLabelError(labelError.TypeUnauthorized, errors.New("unauthorized"))
	errWrongNum := labelError.NewLabelError(labelError.TypeInvalidData, errors.New("invalid data"))
	errCreated := labelError.NewLabelError(labelError.TypeCreated, errors.New("created"))
	errNotFound := labelError.NewLabelError(labelError.TypeNotFound, errors.New("not found"))
	errPaymentRequired := labelError.NewLabelError(labelError.TypeNotEnoughPoints, errors.New("not enough points"))

	mskTz := time.FixedZone("Moscow", 3*3600)

	orderList := []*entity.Order{&entity.Order{
		Number:     entity.OrderNumber("9278923470"),
		Status:     entity.OrderStatusProcessed,
		Accrual:    entity.NewPointValue(500),
		UploadedAt: time.Date(2020, 12, 10, 15, 15, 45, 0, mskTz),
	}, &entity.Order{
		Number:     entity.OrderNumber("12345678903"),
		Status:     entity.OrderStatusProcessing,
		UploadedAt: time.Date(2020, 12, 10, 15, 12, 1, 0, mskTz),
	}, &entity.Order{
		Number:     entity.OrderNumber("346436439"),
		Status:     entity.OrderStatusInvalid,
		UploadedAt: time.Date(2020, 12, 9, 16, 9, 53, 0, mskTz),
	}}

	balance := &entity.UserBalance{
		Current:   entity.NewPointValue(500.5),
		Withdrawn: entity.NewPointValue(42),
	}

	withdrawals := []*entity.Withdrawal{
		&entity.Withdrawal{
			OrderNumber: "2377225624",
			Sum:         entity.NewPointValue(500),
			ProcessedAt: time.Date(2020, 12, 9, 16, 9, 57, 0, mskTz),
		},
	}

	gomock.InOrder(
		mockService.EXPECT().RegisterUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(user, nil).Times(2),
		mockService.EXPECT().RegisterUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(nil, errDef),
		mockService.EXPECT().RegisterUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(nil, errConflict),
	)

	gomock.InOrder(
		mockService.EXPECT().LoginUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(nil, errUnauth),
		mockService.EXPECT().LoginUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(nil, errDef),
		mockService.EXPECT().LoginUser(ctx, gomock.AssignableToTypeOf(login), gomock.AssignableToTypeOf(password)).Return(user, nil),
	)

	gomock.InOrder(
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(errUnauth),
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(errWrongNum),
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(errCreated),
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(errConflict),
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(errDef),
		mockService.EXPECT().AddNewOrder(ctx, gomock.AssignableToTypeOf(sessionID), gomock.AssignableToTypeOf(num)).Return(nil),
	)

	gomock.InOrder(
		mockService.EXPECT().GetOrderList(ctx, gomock.AssignableToTypeOf(sessionID)).Return(orderList, nil),
		mockService.EXPECT().GetOrderList(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errNotFound),
		mockService.EXPECT().GetOrderList(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errUnauth),
		mockService.EXPECT().GetOrderList(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errDef),
	)

	gomock.InOrder(
		mockService.EXPECT().GetUserBalance(ctx, gomock.AssignableToTypeOf(sessionID)).Return(balance, nil),
		mockService.EXPECT().GetUserBalance(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errUnauth),
		mockService.EXPECT().GetUserBalance(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errDef),
	)

	gomock.InOrder(
		mockService.EXPECT().UserOrderWithdraw(
			ctx,
			gomock.AssignableToTypeOf(sessionID),
			gomock.AssignableToTypeOf(num),
			gomock.AssignableToTypeOf(sum),
		).Return(nil),
		mockService.EXPECT().UserOrderWithdraw(
			ctx,
			gomock.AssignableToTypeOf(sessionID),
			gomock.AssignableToTypeOf(num),
			gomock.AssignableToTypeOf(sum),
		).Return(errUnauth),
		mockService.EXPECT().UserOrderWithdraw(
			ctx,
			gomock.AssignableToTypeOf(sessionID),
			gomock.AssignableToTypeOf(num),
			gomock.AssignableToTypeOf(sum),
		).Return(errPaymentRequired),
		mockService.EXPECT().UserOrderWithdraw(
			ctx,
			gomock.AssignableToTypeOf(sessionID),
			gomock.AssignableToTypeOf(num),
			gomock.AssignableToTypeOf(sum),
		).Return(errWrongNum),
		mockService.EXPECT().UserOrderWithdraw(
			ctx,
			gomock.AssignableToTypeOf(sessionID),
			gomock.AssignableToTypeOf(num),
			gomock.AssignableToTypeOf(sum),
		).Return(errDef),
	)

	gomock.InOrder(
		mockService.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(sessionID)).Return(withdrawals, nil),
		mockService.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errNotFound),
		mockService.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errUnauth),
		mockService.EXPECT().GetUserWithdrawals(ctx, gomock.AssignableToTypeOf(sessionID)).Return(nil, errDef),
	)

	return mockService
}
