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

// CheckOtherUserWithDetails mocks base method.
func (m *MockUserRepository) CheckOtherUserWithDetails(arg0 context.Context, arg1 domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOtherUserWithDetails", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckOtherUserWithDetails indicates an expected call of CheckOtherUserWithDetails.
func (mr *MockUserRepositoryMockRecorder) CheckOtherUserWithDetails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOtherUserWithDetails", reflect.TypeOf((*MockUserRepository)(nil).CheckOtherUserWithDetails), arg0, arg1)
}

// DeleteCartItem mocks base method.
func (m *MockUserRepository) DeleteCartItem(arg0 context.Context, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCartItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCartItem indicates an expected call of DeleteCartItem.
func (mr *MockUserRepositoryMockRecorder) DeleteCartItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCartItem", reflect.TypeOf((*MockUserRepository)(nil).DeleteCartItem), arg0, arg1)
}

// FindAddressByID mocks base method.
func (m *MockUserRepository) FindAddressByID(arg0 context.Context, arg1 uint) (domain.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAddressByID", arg0, arg1)
	ret0, _ := ret[0].(domain.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAddressByID indicates an expected call of FindAddressByID.
func (mr *MockUserRepositoryMockRecorder) FindAddressByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAddressByID", reflect.TypeOf((*MockUserRepository)(nil).FindAddressByID), arg0, arg1)
}

// FindAddressByUserID mocks base method.
func (m *MockUserRepository) FindAddressByUserID(arg0 context.Context, arg1 domain.Address, arg2 uint) (domain.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAddressByUserID", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAddressByUserID indicates an expected call of FindAddressByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAddressByUserID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAddressByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAddressByUserID), arg0, arg1, arg2)
}

// FindAllAddressByUserID mocks base method.
func (m *MockUserRepository) FindAllAddressByUserID(arg0 context.Context, arg1 uint) ([]res.ResAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllAddressByUserID", arg0, arg1)
	ret0, _ := ret[0].([]res.ResAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllAddressByUserID indicates an expected call of FindAllAddressByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllAddressByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllAddressByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllAddressByUserID), arg0, arg1)
}

// FindAllCartItemsByCartID mocks base method.
func (m *MockUserRepository) FindAllCartItemsByCartID(arg0 context.Context, arg1 uint) ([]res.ResCartItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllCartItemsByCartID", arg0, arg1)
	ret0, _ := ret[0].([]res.ResCartItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllCartItemsByCartID indicates an expected call of FindAllCartItemsByCartID.
func (mr *MockUserRepositoryMockRecorder) FindAllCartItemsByCartID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllCartItemsByCartID", reflect.TypeOf((*MockUserRepository)(nil).FindAllCartItemsByCartID), arg0, arg1)
}

// FindAllWishListItemsByUserID mocks base method.
func (m *MockUserRepository) FindAllWishListItemsByUserID(arg0 context.Context, arg1 uint) ([]res.ResWishList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllWishListItemsByUserID", arg0, arg1)
	ret0, _ := ret[0].([]res.ResWishList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllWishListItemsByUserID indicates an expected call of FindAllWishListItemsByUserID.
func (mr *MockUserRepositoryMockRecorder) FindAllWishListItemsByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllWishListItemsByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindAllWishListItemsByUserID), arg0, arg1)
}

// FindCartByUserID mocks base method.
func (m *MockUserRepository) FindCartByUserID(arg0 context.Context, arg1 uint) (domain.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCartByUserID", arg0, arg1)
	ret0, _ := ret[0].(domain.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCartByUserID indicates an expected call of FindCartByUserID.
func (mr *MockUserRepositoryMockRecorder) FindCartByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCartByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindCartByUserID), arg0, arg1)
}

// FindCartItemByCartAndProductItemID mocks base method.
func (m *MockUserRepository) FindCartItemByCartAndProductItemID(arg0 context.Context, arg1, arg2 uint) (domain.CartItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCartItemByCartAndProductItemID", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.CartItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCartItemByCartAndProductItemID indicates an expected call of FindCartItemByCartAndProductItemID.
func (mr *MockUserRepositoryMockRecorder) FindCartItemByCartAndProductItemID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCartItemByCartAndProductItemID", reflect.TypeOf((*MockUserRepository)(nil).FindCartItemByCartAndProductItemID), arg0, arg1, arg2)
}

// FindCartItemByID mocks base method.
func (m *MockUserRepository) FindCartItemByID(arg0 context.Context, arg1 uint) (domain.CartItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCartItemByID", arg0, arg1)
	ret0, _ := ret[0].(domain.CartItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCartItemByID indicates an expected call of FindCartItemByID.
func (mr *MockUserRepositoryMockRecorder) FindCartItemByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCartItemByID", reflect.TypeOf((*MockUserRepository)(nil).FindCartItemByID), arg0, arg1)
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

// FindProductItem mocks base method.
func (m *MockUserRepository) FindProductItem(arg0 context.Context, arg1 uint) (domain.ProductItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindProductItem", arg0, arg1)
	ret0, _ := ret[0].(domain.ProductItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindProductItem indicates an expected call of FindProductItem.
func (mr *MockUserRepositoryMockRecorder) FindProductItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindProductItem", reflect.TypeOf((*MockUserRepository)(nil).FindProductItem), arg0, arg1)
}

// FindUser mocks base method.
func (m *MockUserRepository) FindUser(arg0 context.Context, arg1 domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockUserRepositoryMockRecorder) FindUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockUserRepository)(nil).FindUser), arg0, arg1)
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

// SaveCart mocks base method.
func (m *MockUserRepository) SaveCart(arg0 context.Context, arg1 uint) (domain.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCart", arg0, arg1)
	ret0, _ := ret[0].(domain.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveCart indicates an expected call of SaveCart.
func (mr *MockUserRepositoryMockRecorder) SaveCart(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCart", reflect.TypeOf((*MockUserRepository)(nil).SaveCart), arg0, arg1)
}

// SaveCartItem mocks base method.
func (m *MockUserRepository) SaveCartItem(arg0 context.Context, arg1, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCartItem", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCartItem indicates an expected call of SaveCartItem.
func (mr *MockUserRepositoryMockRecorder) SaveCartItem(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCartItem", reflect.TypeOf((*MockUserRepository)(nil).SaveCartItem), arg0, arg1, arg2)
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

// SaveUserWithGoogleDetails mocks base method.
func (m *MockUserRepository) SaveUserWithGoogleDetails(arg0 context.Context, arg1 domain.User) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserWithGoogleDetails", arg0, arg1)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUserWithGoogleDetails indicates an expected call of SaveUserWithGoogleDetails.
func (mr *MockUserRepositoryMockRecorder) SaveUserWithGoogleDetails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserWithGoogleDetails", reflect.TypeOf((*MockUserRepository)(nil).SaveUserWithGoogleDetails), arg0, arg1)
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

// UpdateCartItemQty mocks base method.
func (m *MockUserRepository) UpdateCartItemQty(arg0 context.Context, arg1, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCartItemQty", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCartItemQty indicates an expected call of UpdateCartItemQty.
func (mr *MockUserRepositoryMockRecorder) UpdateCartItemQty(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCartItemQty", reflect.TypeOf((*MockUserRepository)(nil).UpdateCartItemQty), arg0, arg1, arg2)
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