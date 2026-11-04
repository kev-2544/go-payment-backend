package handler

import (
    "context"

    pb "github.com/kev-2544/go-payment-backend/gen/payment"
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

    result, err := h.paymentService.CalculateFee(
        ctx,
        req.Amount,
        req.Currency,
    )

    if err != nil {
        return nil, err
    }

    return &pb.CalculateFeeResponse{
        Amount:      result.Amount,
        Fee:         result.Fee,
        TotalAmount: result.TotalAmount,
    }, nil
}

func (h *PaymentHandler) CreateDeposit(
	ctx contex.Context,
	req *pb.CreateDepositRequest,
) (*pb.CreateDepositResponse, error) {

	tx, err := h.paymentService.CreateDeposit(
		ctx, 
		req.UserId,
		int64(req.Amount),
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDepositResponse{
		Id:       result.ID,
		UserId:   result.UserID,
		Amount:   result.Amount,
		Fee:      result.Fee,
		Currency: result.Currency,
		Status:   result.Status,
	}, nil
}