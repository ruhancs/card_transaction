package factory

import (
	"github.com/ruhancs/card_transaction/domain/infra/gateway"
	"github.com/ruhancs/card_transaction/infra/kafka"
	"github.com/ruhancs/card_transaction/usecase"
)

func TransactionUseCaseFactory(repository gateway.TransactionRepositoryInterface, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	usecase := usecase.UseCaseTransaction{
		TransactionRepository: repository,
		KafkaProducer: producer,
	}

	return usecase
}

