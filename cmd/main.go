package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //driver postgres
	"github.com/ruhancs/card_transaction/infra/grpc/server"
	"github.com/ruhancs/card_transaction/infra/kafka"
	"github.com/ruhancs/card_transaction/infra/repository"
	"github.com/ruhancs/card_transaction/usecase"
	"github.com/ruhancs/card_transaction/usecase/factory"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
	}
	db := setupDB()
	defer db.Close()

	transactionRepository := repository.NewTransactionRepository(db)
	producer := setupKafkaProducer()
	transactionUseCase := factory.TransactionUseCaseFactory(transactionRepository,producer)

	fmt.Println("Server GRPC runing on port 50051...")
	serverGRPC(transactionUseCase)

}

func setupKafkaProducer() kafka.KafkaProducer{
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("BOOTSTRAP_SERVER_KAFKA"))
	return producer
}

func setupDB() *sql.DB {
	db,err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DSN"))
	if err != nil {
		log.Fatal("faied to connecto to db")
	}

	return db
}

func serverGRPC(usecase usecase.UseCaseTransaction) {
	grpcServer := server.NewGrpcServer()
	grpcServer.Usecase = usecase

	grpcServer.Serve()
}