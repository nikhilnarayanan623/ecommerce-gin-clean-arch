// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces (interfaces: AuthUseCase)

// Package mockUseCase is a generated GoMock package.
package mockUseCase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	token "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	interfaces "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	req "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
)

// MockAuthUseCase is a mock of AuthUseCase interface.
type MockAuthUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUseCaseMockRecorder
}

// MockAuthUseCaseMockRecorder is the mock recorder for MockAuthUseCase.
type MockAuthUseCaseMockRecorder struct {
	mock *MockAuthUseCase
}

// NewMockAuthUseCase creates a new mock instance.
func NewMockAuthUseCase(ctrl *gomock.Controller) *MockAuthUseCase {
	mock := &MockAuthUseCase{ctrl: ctrl}
	mock.recorder = &MockAuthUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUseCase) EXPECT() *MockAuthUseCaseMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *MockAuthUseCase) GenerateAccessToken(arg0 context.Context, arg1 interfaces.GenerateTokenParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockAuthUseCaseMockRecorder) GenerateAccessToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockAuthUseCase)(nil).GenerateAccessToken), arg0, arg1)
}

// GenerateRefreshToken mocks base method.
func (m *MockAuthUseCase) GenerateRefreshToken(arg0 context.Context, arg1 interfaces.GenerateTokenParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockAuthUseCaseMockRecorder) GenerateRefreshToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockAuthUseCase)(nil).GenerateRefreshToken), arg0, arg1)
}

// UserLogin mocks base method.
func (m *MockAuthUseCase) UserLogin(arg0 context.Context, arg1 req.Login) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserLogin", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserLogin indicates an expected call of UserLogin.
func (mr *MockAuthUseCaseMockRecorder) UserLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLogin", reflect.TypeOf((*MockAuthUseCase)(nil).UserLogin), arg0, arg1)
}

// VerifyAndGetRefreshTokenSession mocks base method.
func (m *MockAuthUseCase) VerifyAndGetRefreshTokenSession(arg0 context.Context, arg1 string, arg2 token.UserType) (domain.RefreshSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyAndGetRefreshTokenSession", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.RefreshSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyAndGetRefreshTokenSession indicates an expected call of VerifyAndGetRefreshTokenSession.
func (mr *MockAuthUseCaseMockRecorder) VerifyAndGetRefreshTokenSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyAndGetRefreshTokenSession", reflect.TypeOf((*MockAuthUseCase)(nil).VerifyAndGetRefreshTokenSession), arg0, arg1, arg2)
}