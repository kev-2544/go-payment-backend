package provider

import (
	"context"
	"fmt"
)

type FakeProvider struct{}

func NewFakeProvider() *FakeProvider {
	return &FakeProvider{}
}

func (p *FakeProvider) InitiateDeposit(
	ctx context.Context,
	request DepositRequest,
) (DepositResult, error) {

	if err := ctx.Err(); err != nil {
		return DepositResult{}, err
	}

	if request.TransactionID == "" {
		return DepositResult{}, fmt.Errorf("transaction ID is required")
	}

	return DepositResult{
		ProviderReference: "FAKE-MPESA-" + request.TransactionID,
		Status:            "PENDING",
	}, nil
}