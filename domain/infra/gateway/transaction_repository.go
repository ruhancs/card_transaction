package gateway

import "github.com/ruhancs/card_transaction/domain/entity"

type TransactionRepositoryInterface interface {
	SaveTransaction(transaction entity.Transaction, creditCard entity.CreditCard) error
	GetCreditCard(creditCard entity.CreditCard) (entity.CreditCard,error)
	CreateCreditCard(creditCard entity.CreditCard) error
}