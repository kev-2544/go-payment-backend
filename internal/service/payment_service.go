package service

import (
    "context"
	"errors"

	"github.com/google/uuid"
	"github.com/kev-2544/go-payment-backend/internal/model"
	"github.com/kev-2544/go-payment-backend/internal/repository"
)

var (
	ErrInvalidTransactionID = errors.New("transaction ID must be a valid UUID")
	ErrTransactionNotFound  = errors.New("transaction not found")
)

// All amounts are int64 in the currency's smallest unit (e.g. cents).
type FeeResult struct {
    Amount      int64
    Fee         int64
    TotalAmount int64
}

type PaymentService interface {
    CalculateFee(
        ctx context.Context,
        amount int64,
        currency string,
    ) (FeeResult, error)

	CreateDeposit(
		ctx context.Context,
		userId string,
		amount int64,
		currency string,
	)(*model.Transaction, error)

    GetTransaction(
    ctx context.Context,
    id string,
    ) (*model.Transaction, error)
}

type paymentService struct{
	repository repository.PaymentRepository
}



// calculateFee returns a temporary 1% fee, rounded down to the smallest unit.
func calculateFee(amount int64) int64 {
	return amount / 100
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
    amount int64,
    currency string,
) (FeeResult, error) {

    fee := calculateFee(amount)

    return FeeResult{
        Amount:      amount,
        Fee:         fee,
        TotalAmount: amount + fee,
    }, nil
}

func (s *paymentService) GetTransaction(
    ctx context.Context,
    id string,
) (*model.Transaction, error) {

    if _, err := uuid.Parse(id); err != nil {
        return nil, ErrInvalidTransactionID
    }

    tx, err := s.repository.GetTransaction(ctx, id)
    if errors.Is(err, repository.ErrNotFound) {
        return nil, ErrTransactionNotFound
    }
    if err != nil {
        return nil, err
    }

    return tx, nil
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

    fee := calculateFee(amount)

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