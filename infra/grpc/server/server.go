package server

import (
	"log"
	"net"

	"github.com/ruhancs/card_transaction/infra/grpc/pb"
	"github.com/ruhancs/card_transaction/infra/grpc/service"
	"github.com/ruhancs/card_transaction/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	Usecase usecase.UseCaseTransaction
}

func NewGrpcServer() GrpcServer {
	return GrpcServer{}
}

func (s *GrpcServer) Serve() {
	lis,err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatal("could not listen tcp port")
	}

	transactionService := service.NewTransactionService()
	transactionService.ProccessTransactionUseCase = s.Usecase

	//registrar o servico grpc no server
	grpcServer := grpc.NewServer()
	//configuracao para utilizar o grpc server com evans
	reflection.Register(grpcServer)
	pb.RegisterPaymentServiceServer(grpcServer, transactionService)
	
	grpcServer.Serve(lis)
}