// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interfaces/user.go

// Package mockrepo is a generated GoMock package.
package mockrepo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	response "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	domain "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
func (m *MockUserRepository) FindAddressByID(ctx context.Context, addressID uint) (response.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAddressByID", ctx, addressID)
	ret0, _ := ret[0].(response.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAddressByID indicates an expected call of FindAddressByID.
func (mr *MockUserRepositoryMockRecorder) FindAddressByID(ctx, addressID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAddressByID", reflect.TypeOf((*MockUserRepository)(nil).FindAddressByID), ctx, addressID)
}

// FindAllAddressByUserID mocks base method.
func (m *MockUserRepository) FindAllAddressByUserID(ctx context.Context, userID uint) ([]response.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllAddressByUserID", ctx, userID)
	ret0, _ := ret[0].([]response.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllAddressByUserID indicates an expected call of FindAllAddressByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllAddressByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllAddressByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllAddressByUserID), ctx, userID)
}

// FindAllWishListItemsByUserID mocks base method.
func (m *MockUserRepository) FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]response.WishListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllWishListItemsByUserID", ctx, userID)
	ret0, _ := ret[0].([]response.WishListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllWishListItemsByUserID indicates an expected call of FindAllWishListItemsByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllWishListItemsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllWishListItemsByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllWishListItemsByUserID), ctx, userID)
}

// FindCountryByID mocks base method.
func (m *MockUserRepository) FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCountryByID", ctx, countryID)
	ret0, _ := ret[0].(domain.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCountryByID indicates an expected call of FindCountryByID.
func (mr *MockUserRepositoryMockRecorder) FindCountryByID(ctx, countryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCountryByID", reflect.TypeOf((*MockUserRepository)(nil).FindCountryByID), ctx, countryID)
}

// FindUserByEmail mocks base method.
func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockUserRepositoryMockRecorder) FindUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindUserByEmail), ctx, email)
}

// FindUserByPhoneNumber mocks base method.
func (m *MockUserRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByPhoneNumber", ctx, phoneNumber)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByPhoneNumber indicates an expected call of FindUserByPhoneNumber.
func (mr *MockUserRepositoryMockRecorder) FindUserByPhoneNumber(ctx, phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByPhoneNumber", reflect.TypeOf((*MockUserRepository)(nil).FindUserByPhoneNumber), ctx, phoneNumber)
}

// FindUserByUserID mocks base method.
func (m *MockUserRepository) FindUserByUserID(ctx context.Context, userID uint) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserID", ctx, userID)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserID indicates an expected call of FindUserByUserID.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserID), ctx, userID)
}

// FindUserByUserName mocks base method.
func (m *MockUserRepository) FindUserByUserName(ctx context.Context, userName string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserName", ctx, userName)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserName indicates an expected call of FindUserByUserName.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserName(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserName", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserName), ctx, userName)
}

// FindUserByUserNameEmailOrPhoneNotID mocks base method.
func (m *MockUserRepository) FindUserByUserNameEmailOrPhoneNotID(ctx context.Context, user domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByUserNameEmailOrPhoneNotID", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByUserNameEmailOrPhoneNotID indicates an expected call of FindUserByUserNameEmailOrPhoneNotID.
func (mr *MockUserRepositoryMockRecorder) FindUserByUserNameEmailOrPhoneNotID(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByUserNameEmailOrPhoneNotID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByUserNameEmailOrPhoneNotID), ctx, user)
}

// FindWishListItem mocks base method.
func (m *MockUserRepository) FindWishListItem(ctx context.Context, productID, userID uint) (domain.WishList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindWishListItem", ctx, productID, userID)
	ret0, _ := ret[0].(domain.WishList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWishListItem indicates an expected call of FindWishListItem.
func (mr *MockUserRepositoryMockRecorder) FindWishListItem(ctx, productID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWishListItem", reflect.TypeOf((*MockUserRepository)(nil).FindWishListItem), ctx, productID, userID)
}

// IsAddressAlreadyExistForUser mocks base method.
func (m *MockUserRepository) IsAddressAlreadyExistForUser(ctx context.Context, address domain.Address, userID uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAddressAlreadyExistForUser", ctx, address, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAddressAlreadyExistForUser indicates an expected call of IsAddressAlreadyExistForUser.
func (mr *MockUserRepositoryMockRecorder) IsAddressAlreadyExistForUser(ctx, address, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAddressAlreadyExistForUser", reflect.TypeOf((*MockUserRepository)(nil).IsAddressAlreadyExistForUser), ctx, address, userID)
}

// IsAddressIDExist mocks base method.
func (m *MockUserRepository) IsAddressIDExist(ctx context.Context, addressID uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAddressIDExist", ctx, addressID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAddressIDExist indicates an expected call of IsAddressIDExist.
func (mr *MockUserRepositoryMockRecorder) IsAddressIDExist(ctx, addressID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAddressIDExist", reflect.TypeOf((*MockUserRepository)(nil).IsAddressIDExist), ctx, addressID)
}

// RemoveWishListItem mocks base method.
func (m *MockUserRepository) RemoveWishListItem(ctx context.Context, userID, productItemID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveWishListItem", ctx, userID, productItemID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveWishListItem indicates an expected call of RemoveWishListItem.
func (mr *MockUserRepositoryMockRecorder) RemoveWishListItem(ctx, userID, productItemID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveWishListItem", reflect.TypeOf((*MockUserRepository)(nil).RemoveWishListItem), ctx, userID, productItemID)
}

// SaveAddress mocks base method.
func (m *MockUserRepository) SaveAddress(ctx context.Context, address domain.Address) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAddress", ctx, address)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveAddress indicates an expected call of SaveAddress.
func (mr *MockUserRepositoryMockRecorder) SaveAddress(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAddress", reflect.TypeOf((*MockUserRepository)(nil).SaveAddress), ctx, address)
}

// SaveUser mocks base method.
func (m *MockUserRepository) SaveUser(ctx context.Context, user domain.User) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", ctx, user)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockUserRepositoryMockRecorder) SaveUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockUserRepository)(nil).SaveUser), ctx, user)
}

// SaveUserAddress mocks base method.
func (m *MockUserRepository) SaveUserAddress(ctx context.Context, userAdress domain.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserAddress", ctx, userAdress)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUserAddress indicates an expected call of SaveUserAddress.
func (mr *MockUserRepositoryMockRecorder) SaveUserAddress(ctx, userAdress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserAddress", reflect.TypeOf((*MockUserRepository)(nil).SaveUserAddress), ctx, userAdress)
}

// SaveWishListItem mocks base method.
func (m *MockUserRepository) SaveWishListItem(ctx context.Context, wishList domain.WishList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWishListItem", ctx, wishList)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWishListItem indicates an expected call of SaveWishListItem.
func (mr *MockUserRepositoryMockRecorder) SaveWishListItem(ctx, wishList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWishListItem", reflect.TypeOf((*MockUserRepository)(nil).SaveWishListItem), ctx, wishList)
}

// UpdateAddress mocks base method.
func (m *MockUserRepository) UpdateAddress(ctx context.Context, address domain.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAddress", ctx, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAddress indicates an expected call of UpdateAddress.
func (mr *MockUserRepositoryMockRecorder) UpdateAddress(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAddress", reflect.TypeOf((*MockUserRepository)(nil).UpdateAddress), ctx, address)
}

// UpdateBlockStatus mocks base method.
func (m *MockUserRepository) UpdateBlockStatus(ctx context.Context, userID uint, blockStatus bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBlockStatus", ctx, userID, blockStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBlockStatus indicates an expected call of UpdateBlockStatus.
func (mr *MockUserRepositoryMockRecorder) UpdateBlockStatus(ctx, userID, blockStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBlockStatus", reflect.TypeOf((*MockUserRepository)(nil).UpdateBlockStatus), ctx, userID, blockStatus)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), ctx, user)
}

// UpdateUserAddress mocks base method.
func (m *MockUserRepository) UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserAddress", ctx, userAddress)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserAddress indicates an expected call of UpdateUserAddress.
func (mr *MockUserRepositoryMockRecorder) UpdateUserAddress(ctx, userAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserAddress", reflect.TypeOf((*MockUserRepository)(nil).UpdateUserAddress), ctx, userAddress)
}

// UpdateVerified mocks base method.
func (m *MockUserRepository) UpdateVerified(ctx context.Context, userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVerified", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVerified indicates an expected call of UpdateVerified.
func (mr *MockUserRepositoryMockRecorder) UpdateVerified(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVerified", reflect.TypeOf((*MockUserRepository)(nil).UpdateVerified), ctx, userID)
}
