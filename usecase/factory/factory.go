package factory

import (
	"github.com/ruhancs/card_transaction/domain/infra/gateway"
	"github.com/ruhancs/card_transaction/usecase"
)

func TransactionUseCaseFactory(repository gateway.TransactionRepositoryInterface) usecase.UseCaseTransaction {
	usecase := usecase.UseCaseTransaction{
		TransactionRepository: repository,
	}

	return usecase
}