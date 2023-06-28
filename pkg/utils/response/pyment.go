package response

import "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"

type OrderPayment struct {
	PaymentType  domain.PaymentType `json:"payment_type"`
	PaymentOrder any                `json:"payment_order"`
}
