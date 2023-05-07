// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces (interfaces: AuthRepository)

// Package mockRepository is a generated GoMock package.
package mockRepository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

// MockAuthRepository is a mock of AuthRepository interface.
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository.
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance.
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// FindRefreshSessionByTokenID mocks base method.
func (m *MockAuthRepository) FindRefreshSessionByTokenID(arg0 context.Context, arg1 uuid.UUID) (domain.RefreshSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRefreshSessionByTokenID", arg0, arg1)
	ret0, _ := ret[0].(domain.RefreshSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRefreshSessionByTokenID indicates an expected call of FindRefreshSessionByTokenID.
func (mr *MockAuthRepositoryMockRecorder) FindRefreshSessionByTokenID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRefreshSessionByTokenID", reflect.TypeOf((*MockAuthRepository)(nil).FindRefreshSessionByTokenID), arg0, arg1)
}

// SaveRefreshSession mocks base method.
func (m *MockAuthRepository) SaveRefreshSession(arg0 context.Context, arg1 domain.RefreshSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRefreshSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRefreshSession indicates an expected call of SaveRefreshSession.
func (mr *MockAuthRepositoryMockRecorder) SaveRefreshSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRefreshSession", reflect.TypeOf((*MockAuthRepository)(nil).SaveRefreshSession), arg0, arg1)
}