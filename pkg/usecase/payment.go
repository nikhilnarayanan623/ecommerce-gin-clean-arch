package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
)

type paymentUseCase struct {
	paymentRepo interfaces.PaymentRepository
}

func NewPaymentUseCase(paymentRepo interfaces.PaymentRepository) service.PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
	}
}

func (c *paymentUseCase) FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error) {
	return c.paymentRepo.FindAllPaymentMethods(ctx)
}

func (c *paymentUseCase) FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error) {
	return c.paymentRepo.FindPaymentMethodByID(ctx, paymentMethodID)
}

func (c *paymentUseCase) SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error {

	// first check the payment_method alreadcy exist with given payment_type
	checkPayment, err := c.paymentRepo.FindPaymentMethodByType(ctx, paymentMethod.PaymentType)
	if err != nil {
		return err
	} else if checkPayment.ID != 0 {
		return fmt.Errorf("an payment_method already exist wtih given payment_type %v", paymentMethod.PaymentType)
	}

	// save payment
	paymentMethodID, err := c.paymentRepo.SavePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return err
	}

	log.Printf("successfully saved payment method for payment_type %v with id %v", paymentMethod.PaymentType, paymentMethodID)
	return nil
}
func (c *paymentUseCase) UpdatePaymentMethod(ctx context.Context, paymentMethod request.PaymentMethodUpdate) error {

	// first check the given payement_method_id is valid or not
	checkPayment, err := c.paymentRepo.FindPaymentMethodByID(ctx, paymentMethod.ID)
	if err != nil {
		return err
	} else if checkPayment.ID == 0 {
		return fmt.Errorf("invalid payment_method_id %v", paymentMethod.ID)
	}

	// // check the given payment_type already exist
	// checkPayment, err = c.paymentRepo.FindPaymentMethodByType(ctx, paymentMethod.PaymentType)
	// if err != nil {
	// 	return err
	// } else if checkPayment.ID != 0 && checkPayment.ID != paymentMethod.ID {
	// 	return fmt.Errorf("an payment_method already exist wtih given payment_type %v", paymentMethod.PaymentType)
	// }

	err = c.paymentRepo.UpdatePaymentMethod(ctx, paymentMethod)
	if err != nil {
		return err
	}

	log.Printf("successfully update payment_method of payment_method_id %v ", paymentMethod.ID)
	return nil
}
