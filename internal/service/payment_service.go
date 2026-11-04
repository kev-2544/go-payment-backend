package service

import (
    "context"
	"errors"
	"github.com/kev-2544/go-payment-backend/internal/model"

	"github.com/kev-2544/go-payment-backend/internal/repository"
	"github.com/kev-2544/go-payment-backend/internal/model"

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

	CreateDeposit(
		ctx context.Context,
		userId string,
		amount int64,
		currency string,
	)(*model.Transaction, error)
}

type paymentService struct{
	repository repository.PaymentRepository
}



func NewPaymentService(
	repo repository.PaymentRepository,
) PaymentService {
    return &paymentService{
		repository: repo,
	}
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

func (s *paymentService) CreateDeposit(
    ctx context.Context,
    userID string,
    amount int64,
    currency string,
) (*model.Transaction, error) {

    if amount <= 0 {
        return nil, errors.New("amount must be greater than zero")
    }

    if userID == "" {
        return nil, errors.New("user ID is required")
    }

    fee := amount / 100 // temporary 1% fee

    tx := &model.Transaction{
        ID:       uuid.NewString(),
        UserID:   userID,
        Amount:   amount,
        Fee:      fee,
        Currency: currency,
        Status:   "PENDING",
    }

    err := s.repository.CreateTransaction(ctx, tx)
    if err != nil {
        return nil, err
    }

    return tx, nil
}