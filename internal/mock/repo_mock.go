// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/AnnV0lokitina/diplom/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockRepo) AddOrder(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, user, orderNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockRepoMockRecorder) AddOrder(ctx, user, orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockRepo)(nil).AddOrder), ctx, user, orderNumber)
}

// AddOrderInfo mocks base method.
func (m *MockRepo) AddOrderInfo(ctx context.Context, orderNumber entity.OrderNumber, status entity.OrderStatus, accrual entity.PointValue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrderInfo", ctx, orderNumber, status, accrual)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrderInfo indicates an expected call of AddOrderInfo.
func (mr *MockRepoMockRecorder) AddOrderInfo(ctx, orderNumber, status, accrual interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrderInfo", reflect.TypeOf((*MockRepo)(nil).AddOrderInfo), ctx, orderNumber, status, accrual)
}

// AddUserSession mocks base method.
func (m *MockRepo) AddUserSession(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserSession", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserSession indicates an expected call of AddUserSession.
func (mr *MockRepoMockRecorder) AddUserSession(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserSession", reflect.TypeOf((*MockRepo)(nil).AddUserSession), ctx, user)
}

// AuthUser mocks base method.
func (m *MockRepo) AuthUser(ctx context.Context, login, passwordHash string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUser", ctx, login, passwordHash)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthUser indicates an expected call of AuthUser.
func (mr *MockRepoMockRecorder) AuthUser(ctx, login, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUser", reflect.TypeOf((*MockRepo)(nil).AuthUser), ctx, login, passwordHash)
}

// Close mocks base method.
func (m *MockRepo) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepoMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepo)(nil).Close), ctx)
}

// CreateUser mocks base method.
func (m *MockRepo) CreateUser(ctx context.Context, sessionID, login, passwordHash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, sessionID, login, passwordHash)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepoMockRecorder) CreateUser(ctx, sessionID, login, passwordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepo)(nil).CreateUser), ctx, sessionID, login, passwordHash)
}

// GetOrdersListForRequest mocks base method.
func (m *MockRepo) GetOrdersListForRequest(ctx context.Context) ([]entity.OrderNumber, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersListForRequest", ctx)
	ret0, _ := ret[0].([]entity.OrderNumber)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersListForRequest indicates an expected call of GetOrdersListForRequest.
func (mr *MockRepoMockRecorder) GetOrdersListForRequest(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersListForRequest", reflect.TypeOf((*MockRepo)(nil).GetOrdersListForRequest), ctx)
}

// GetUserBalance mocks base method.
func (m *MockRepo) GetUserBalance(ctx context.Context, user *entity.User) (*entity.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", ctx, user)
	ret0, _ := ret[0].(*entity.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockRepoMockRecorder) GetUserBalance(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockRepo)(nil).GetUserBalance), ctx, user)
}

// GetUserBySessionID mocks base method.
func (m *MockRepo) GetUserBySessionID(ctx context.Context, activeSessionID string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBySessionID", ctx, activeSessionID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBySessionID indicates an expected call of GetUserBySessionID.
func (mr *MockRepoMockRecorder) GetUserBySessionID(ctx, activeSessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBySessionID", reflect.TypeOf((*MockRepo)(nil).GetUserBySessionID), ctx, activeSessionID)
}

// GetUserOrders mocks base method.
func (m *MockRepo) GetUserOrders(ctx context.Context, user *entity.User) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", ctx, user)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockRepoMockRecorder) GetUserOrders(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockRepo)(nil).GetUserOrders), ctx, user)
}

// GetUserWithdrawals mocks base method.
func (m *MockRepo) GetUserWithdrawals(ctx context.Context, user *entity.User) ([]*entity.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithdrawals", ctx, user)
	ret0, _ := ret[0].([]*entity.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithdrawals indicates an expected call of GetUserWithdrawals.
func (mr *MockRepoMockRecorder) GetUserWithdrawals(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithdrawals", reflect.TypeOf((*MockRepo)(nil).GetUserWithdrawals), ctx, user)
}

// UserOrderWithdraw mocks base method.
func (m *MockRepo) UserOrderWithdraw(ctx context.Context, user *entity.User, orderNumber entity.OrderNumber, sum entity.PointValue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOrderWithdraw", ctx, user, orderNumber, sum)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserOrderWithdraw indicates an expected call of UserOrderWithdraw.
func (mr *MockRepoMockRecorder) UserOrderWithdraw(ctx, user, orderNumber, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOrderWithdraw", reflect.TypeOf((*MockRepo)(nil).UserOrderWithdraw), ctx, user, orderNumber, sum)
}