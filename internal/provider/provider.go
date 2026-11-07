package provider

import "context"

type DepositRequest struct {
    TransactionID string
    UserID        string
    Amount        int64
    Currency      string
}

type DepositResult struct {
    ProviderReference string
    Status            string
}

type PaymentProvider interface {
    InitiateDeposit(
        ctx context.Context,
        request DepositRequest,
    ) (DepositResult, error)
}