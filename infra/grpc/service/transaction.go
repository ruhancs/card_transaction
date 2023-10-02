package service

import (
	"context"

	"github.com/ruhancs/card_transaction/dto"
	"github.com/ruhancs/card_transaction/infra/grpc/pb"
	"github.com/ruhancs/card_transaction/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionService struct {
	ProccessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService()*TransactionService {
	return &TransactionService{}
}

func (service TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*emptypb.Empty, error) {
	//inserir o PaymentRequest do grpc em TransactionDTO para utilizar o usecase
	transactionDTO := dto.Transaction{
		Name: in.GetCreditCard().GetName(),
		Number: in.GetCreditCard().GetNumber(),
		ExpirationMonth: in.GetCreditCard().GetExpirationMonth(),
		ExpirationYear: in.GetCreditCard().GetExpirationYear(),
		Amount: in.GetAmount(),
		CVV: in.CreditCard.GetCvv(),
		Store: in.GetStore(),
		Description: in.GetDescription(),
	}

	transaction,err := service.ProccessTransactionUseCase.ProccessTransaction(transactionDTO)
	if err != nil {
		return &emptypb.Empty{},status.Error(codes.FailedPrecondition, err.Error())
	}
	if transaction.Status != "approved" {
		return &emptypb.Empty{},status.Error(codes.FailedPrecondition, "transaction reject, please enter in contact with your bank")
	}

	return &emptypb.Empty{},nil
}