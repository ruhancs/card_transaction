package usecase

import (
	"encoding/json"
	"time"

	"github.com/ruhancs/card_transaction/domain/entity"
	"github.com/ruhancs/card_transaction/domain/infra/gateway"
	"github.com/ruhancs/card_transaction/dto"
	"github.com/ruhancs/card_transaction/infra/kafka"
)

type UseCaseTransaction struct {
	TransactionRepository gateway.TransactionRepositoryInterface
	KafkaProducer kafka.KafkaProducer
}

func NewTransactionUseCase(repository gateway.TransactionRepositoryInterface) UseCaseTransaction {
	return UseCaseTransaction{
		TransactionRepository: repository,
	}
}

func (usecase *UseCaseTransaction) ProccessTransaction(transactionDTO dto.Transaction) (entity.Transaction,error) {
	creditCard := entity.NewCreditCard()
	creditCard.Name = transactionDTO.Name
	creditCard.Number = transactionDTO.Number
	creditCard.ExpirationMonth = transactionDTO.ExpirationMonth
	creditCard.ExpirationYear = transactionDTO.ExpirationYear
	creditCard.CVV = transactionDTO.CVV

	creditCardBalanceAndLimit,err := usecase.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return entity.Transaction{},err
	}
	creditCard.ID = creditCardBalanceAndLimit.ID
	creditCard.Limit = creditCardBalanceAndLimit.Limit
	creditCard.Balance = creditCardBalanceAndLimit.Balance

	t := entity.NewTransaction()
	t.CreditCardID = creditCard.ID
	t.Amount = transactionDTO.Amount
	t.Status = transactionDTO.Store
	t.Description = transactionDTO.Description
	t.CreatedAt = time.Now()

	//validacao do limit para transacao
	t.ProccessAndValidate(creditCard)

	err = usecase.TransactionRepository.SaveTransaction(*t,*creditCard)
	if err != nil {
		return entity.Transaction{},err
	}

	transactionDTO.ID = t.ID
	transactionDTO.CreatedAt = t.CreatedAt
	//transformar transactionDTO em json
	transactionJson,err := json.Marshal(transactionDTO)
	if err != nil {
		return entity.Transaction{},err
	}
	
	//publicar a msg no topico
	err = usecase.KafkaProducer.Publish(string(transactionJson), "payments")
	if err != nil {
		return entity.Transaction{},err
	}


	return *t,nil
}