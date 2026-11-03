package service

import (
    "context"
)

type FeeResult struct {
    Amount      float64
    Fee         float64
    TotalAmount float64
}

type PaymentService interface {
    CalculateFee(
        ctx context.Context,
        amount float64,
        currency string,
    ) (FeeResult, error)
}

type paymentService struct{}

func NewPaymentService() PaymentService {
    return &paymentService{}
}

func (s *paymentService) CalculateFee(
    ctx context.Context,
    amount float64,
    currency string,
) (FeeResult, error) {

    fee := amount * 0.01

    return FeeResult{
        Amount:      amount,
        Fee:         fee,
        TotalAmount: amount + fee,
    }, nil
}