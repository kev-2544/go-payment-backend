package repository

import (
    "context"
    "errors"

    "github.com/kev-2544/go-payment-backend/internal/model"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("not found")

type PaymentRepository interface {
    CreateTransaction(
        ctx context.Context,
        tx *model.Transaction,
    ) error

    GetTransaction(
        ctx context.Context,
        id string,
    ) (*model.Transaction, error)
}