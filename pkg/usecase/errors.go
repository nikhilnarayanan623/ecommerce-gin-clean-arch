package usecase

import "errors"

var (
	// login
	ErrEmptyLoginCredentials = errors.New("all login credentials are empty")
	ErrUserNotExist          = errors.New("user not exist with given login credentials")
	ErrUserBlocked           = errors.New("user blocked by admin")
	ErrWrongPassword         = errors.New("password doesn't match")
	// otp
	ErrInvalidOtpID = errors.New("invalid otp id")
	ErrOtpExpired   = errors.New("otp session expired")
	ErrInvalidOtp   = errors.New("invalid otp")

	// refresh token
	ErrRefreshSessionNotExist = errors.New("there is no refresh token session for this token")
	ErrRefreshSessionExpired  = errors.New("refresh token expired in session")
	ErrRefreshSessionBlocked  = errors.New("refresh token blocked in session")

	// signup
	ErrUserAlreadyExit = errors.New("user already exist")

	// cart
	ErrInvalidProductItemID  = errors.New("invalid product_item_id")
	ErrProductItemOutOfStock = errors.New("product is now out of stock")
	ErrCartItemAlreadyExist  = errors.New("product_item already exist on the cart")
	ErrCartItemNotExit       = errors.New("product_item not exist on cart")
	ErrEmptyCart             = errors.New("user cart is empty")

	ErrRequireMinimumCartItemQty = errors.New("update cart item qty can not less than 1")
	ErrInvalidCartItemUpdateQty  = errors.New("update cart item qty reached max limit")

	// admin
	ErrSameBlockStatus = errors.New("user block status already in given status")
)
