package repository

import (
    "context"

    "github.com/kev-2544/go-payment-backend/internal/model"
)

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