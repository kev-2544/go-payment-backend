package main

import (
	"context"
	"log"
	"net"

	pb "github.com/kev-2544/go-payment-backend/gen/payment"
	"github.com/kev-2544/go-payment-backend/internal/handler"
	"github.com/kev-2544/go-payment-backend/internal/repository"
	"github.com/kev-2544/go-payment-backend/internal/service"

	"google.golang.org/grpc"
)

func main() {
	// 1. Create a context
	ctx := context.Background()

	// 2. Connect to PostgreSQL
	dbPool, err := repository.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 3. Create the repository
	paymentRepo := repository.NewPaymentRepository(dbPool)

	// 4. Create the service
	paymentSvc := service.NewPaymentService(paymentRepo)

	// 5. Create the handler
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	// 6. Create the gRPC server
	grpcServer := grpc.NewServer()

	// 7. Register our payment handler
	pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// 8. Open a network listener
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server listening on :50051")

	// 9. Start accepting gRPC requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}