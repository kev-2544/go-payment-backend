package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/kev-2544/go-payment-backend/gen/payment"
	"github.com/kev-2544/go-payment-backend/internal/money"
	"github.com/kev-2544/go-payment-backend/internal/service"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentHandler(
	paymentService service.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CalculateFee(
	ctx context.Context,
	req *pb.CalculateFeeRequest,
) (*pb.CalculateFeeResponse, error) {

	currency, amount, err := parseAmount(req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	result, err := h.paymentService.CalculateFee(ctx, amount, currency)
	if err != nil {
		return nil, err
	}

	return &pb.CalculateFeeResponse{
		Amount:      formatAmount(result.Amount, currency),
		Fee:         formatAmount(result.Fee, currency),
		TotalAmount: formatAmount(result.TotalAmount, currency),
	}, nil
}

func (h *PaymentHandler) CreateDeposit(
	ctx context.Context,
	req *pb.CreateDepositRequest,
) (*pb.CreateDepositResponse, error) {

	currency, amount, err := parseAmount(req.Amount, req.Currency)
	if err != nil {
		return nil, err
	}

	tx, err := h.paymentService.CreateDeposit(ctx, req.UserId, amount, currency)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDepositResponse{
		TransactionId: tx.ID,
		Status:        tx.Status,
		Amount:        formatAmount(tx.Amount, tx.Currency),
		Fee:           formatAmount(tx.Fee, tx.Currency),
	}, nil
}

func (h *PaymentHandler) GetTransaction(
	ctx context.Context,
	req *pb.GetTransactionRequest,
) (*pb.GetTransactionResponse, error) {

	tx, err := h.paymentService.GetTransaction(ctx, req.TransactionId)
	switch {
	case errors.Is(err, service.ErrInvalidTransactionID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrTransactionNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, "failed to get transaction")
	}

	return &pb.GetTransactionResponse{
		TransactionId: tx.ID,
		Status:        tx.Status,
		Amount:        formatAmount(tx.Amount, tx.Currency),
		Fee:           formatAmount(tx.Fee, tx.Currency),
		Currency:      tx.Currency,
	}, nil
}

// parseAmount validates the currency and converts a decimal amount string
// (e.g. "10.50") into minor units (1050).
func parseAmount(amount, currency string) (string, int64, error) {
	code, err := money.NormalizeCurrency(currency)
	if err != nil {
		return "", 0, status.Error(codes.InvalidArgument, err.Error())
	}
	minor, err := money.Parse(amount, code)
	if err != nil {
		return "", 0, status.Error(codes.InvalidArgument, err.Error())
	}
	return code, minor, nil
}

// formatAmount converts minor units back to a decimal string. The currency
// was validated on the way in, so Format cannot fail here.
func formatAmount(minor int64, currency string) string {
	s, _ := money.Format(minor, currency)
	return s
}
