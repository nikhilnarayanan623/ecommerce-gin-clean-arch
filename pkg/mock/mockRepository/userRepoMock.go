// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces (interfaces: UserRepository)

// Package mockRepository is a generated GoMock package.
package mockRepository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	res "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindAddressByID mocks base method.
func (m *MockUserRepository) FindAddressByID(arg0 context.Context, arg1 uint) (res.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAddressByID", arg0, arg1)
	ret0, _ := ret[0].(res.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAddressByID indicates an expected call of FindAddressByID.
func (mr *MockUserRepositoryMockRecorder) FindAddressByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAddressByID", reflect.TypeOf((*MockUserRepository)(nil).FindAddressByID), arg0, arg1)
}

// FindAllAddressByUserID mocks base method.
func (m *MockUserRepository) FindAllAddressByUserID(arg0 context.Context, arg1 uint) ([]res.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllAddressByUserID", arg0, arg1)
	ret0, _ := ret[0].([]res.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllAddressByUserID indicates an expected call of FindAllAddressByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllAddressByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllAddressByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllAddressByUserID), arg0, arg1)
}

// FindAllWishListItemsByUserID mocks base method.
func (m *MockUserRepository) FindAllWishListItemsByUserID(arg0 context.Context, arg1 uint) ([]res.WishList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllWishListItemsByUserID", arg0, arg1)
	ret0, _ := ret[0].([]res.WishList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllWishListItemsByUserID indicates an expected call of FindAllWishListItemsByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllWishListItemsByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllWishListItemsByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllWishListItemsByUserID), arg0, arg1)
}

// FindCountryByID mocks base method.
func (m *MockUserRepository) FindCountryByID(arg0 context.Context, arg1 uint) (domain.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCountryByID", arg0, arg1)
	ret0, _ := ret[0].(domain.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCountryByID indicates an expected call of FindCountryByID.
func (mr *MockUserRepositoryMockRecorder) FindCountryByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCountryByID", reflect.TypeOf((*MockUserRepository)(nil).FindCountryByID), arg0, arg1)
}

// FindUserByEmail mocks base method.
func (m *MockUserRepository) FindUserByEmail(arg0 context.Context, arg1 string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockUserRepositoryMockRecorder) FindUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindUserByEmail), arg0, arg1)
}

// FindUserByPhoneNumber mocks base method.
func (m *MockUserRepository) FindUserByPhoneNumber(arg0 context.Context, arg1 string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByPhoneNumber", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByPhoneNumber indicates an expected call of FindUserByPhoneNumber.
func (mr *MockUserRepositoryMockRecorder) FindUserByPhoneNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByPhoneNumber", reflect.TypeOf((*MockUserRepository)(nil).FindUserByPhoneNumber), arg0, arg1)
}

// FindUserByUserID mocks base method.
func (m *MockUserRepository) FindUserByUserID(arg0 context.Context, arg1 uint) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserID", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserID indicates an expected call of FindUserByUserID.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserID), arg0, arg1)
}

// FindUserByUserName mocks base method.
func (m *MockUserRepository) FindUserByUserName(arg0 context.Context, arg1 string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserName", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserName indicates an expected call of FindUserByUserName.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserName", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserName), arg0, arg1)
}

// FindUserByUserNameEmailOrPhoneNotID mocks base method.
func (m *MockUserRepository) FindUserByUserNameEmailOrPhoneNotID(arg0 context.Context, arg1 domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserNameEmailOrPhoneNotID", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserNameEmailOrPhoneNotID indicates an expected call of FindUserByUserNameEmailOrPhoneNotID.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserNameEmailOrPhoneNotID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserNameEmailOrPhoneNotID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserNameEmailOrPhoneNotID), arg0, arg1)
}

// FindWishListItem mocks base method.
func (m *MockUserRepository) FindWishListItem(arg0 context.Context, arg1, arg2 uint) (domain.WishList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindWishListItem", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.WishList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWishListItem indicates an expected call of FindWishListItem.
func (mr *MockUserRepositoryMockRecorder) FindWishListItem(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWishListItem", reflect.TypeOf((*MockUserRepository)(nil).FindWishListItem), arg0, arg1, arg2)
}

// IsAddressAlreadyExistForUser mocks base method.
func (m *MockUserRepository) IsAddressAlreadyExistForUser(arg0 context.Context, arg1 domain.Address, arg2 uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAddressAlreadyExistForUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAddressAlreadyExistForUser indicates an expected call of IsAddressAlreadyExistForUser.
func (mr *MockUserRepositoryMockRecorder) IsAddressAlreadyExistForUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAddressAlreadyExistForUser", reflect.TypeOf((*MockUserRepository)(nil).IsAddressAlreadyExistForUser), arg0, arg1, arg2)
}

// IsAddressIDExist mocks base method.
func (m *MockUserRepository) IsAddressIDExist(arg0 context.Context, arg1 uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAddressIDExist", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAddressIDExist indicates an expected call of IsAddressIDExist.
func (mr *MockUserRepositoryMockRecorder) IsAddressIDExist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAddressIDExist", reflect.TypeOf((*MockUserRepository)(nil).IsAddressIDExist), arg0, arg1)
}

// RemoveWishListItem mocks base method.
func (m *MockUserRepository) RemoveWishListItem(arg0 context.Context, arg1 domain.WishList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveWishListItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveWishListItem indicates an expected call of RemoveWishListItem.
func (mr *MockUserRepositoryMockRecorder) RemoveWishListItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveWishListItem", reflect.TypeOf((*MockUserRepository)(nil).RemoveWishListItem), arg0, arg1)
}

// SaveAddress mocks base method.
func (m *MockUserRepository) SaveAddress(arg0 context.Context, arg1 domain.Address) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAddress", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveAddress indicates an expected call of SaveAddress.
func (mr *MockUserRepositoryMockRecorder) SaveAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAddress", reflect.TypeOf((*MockUserRepository)(nil).SaveAddress), arg0, arg1)
}

// SaveUser mocks base method.
func (m *MockUserRepository) SaveUser(arg0 context.Context, arg1 domain.User) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockUserRepositoryMockRecorder) SaveUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockUserRepository)(nil).SaveUser), arg0, arg1)
}

// SaveUserAddress mocks base method.
func (m *MockUserRepository) SaveUserAddress(arg0 context.Context, arg1 domain.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserAddress", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUserAddress indicates an expected call of SaveUserAddress.
func (mr *MockUserRepositoryMockRecorder) SaveUserAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserAddress", reflect.TypeOf((*MockUserRepository)(nil).SaveUserAddress), arg0, arg1)
}

// SaveWishListItem mocks base method.
func (m *MockUserRepository) SaveWishListItem(arg0 context.Context, arg1 domain.WishList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWishListItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWishListItem indicates an expected call of SaveWishListItem.
func (mr *MockUserRepositoryMockRecorder) SaveWishListItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWishListItem", reflect.TypeOf((*MockUserRepository)(nil).SaveWishListItem), arg0, arg1)
}

// UpdateAddress mocks base method.
func (m *MockUserRepository) UpdateAddress(arg0 context.Context, arg1 domain.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAddress", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAddress indicates an expected call of UpdateAddress.
func (mr *MockUserRepositoryMockRecorder) UpdateAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAddress", reflect.TypeOf((*MockUserRepository)(nil).UpdateAddress), arg0, arg1)
}

// UpdateBlockStatus mocks base method.
func (m *MockUserRepository) UpdateBlockStatus(arg0 context.Context, arg1 uint, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBlockStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBlockStatus indicates an expected call of UpdateBlockStatus.
func (mr *MockUserRepositoryMockRecorder) UpdateBlockStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBlockStatus", reflect.TypeOf((*MockUserRepository)(nil).UpdateBlockStatus), arg0, arg1, arg2)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(arg0 context.Context, arg1 domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), arg0, arg1)
}

// UpdateUserAddress mocks base method.
func (m *MockUserRepository) UpdateUserAddress(arg0 context.Context, arg1 domain.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserAddress", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserAddress indicates an expected call of UpdateUserAddress.
func (mr *MockUserRepositoryMockRecorder) UpdateUserAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserAddress", reflect.TypeOf((*MockUserRepository)(nil).UpdateUserAddress), arg0, arg1)
}
