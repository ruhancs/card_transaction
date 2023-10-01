package repository

import (
	"database/sql"
	"errors"

	"github.com/ruhancs/card_transaction/domain/entity"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (repository *TransactionRepository) SaveTransaction(transaction entity.Transaction, creditCard entity.CreditCard) error {
	stmt,err := repository.DB.Prepare(`insert into transactions(id, credit_card_id, amount, status, description, store, created_at)
		values($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		return err
	}
	_,err = stmt.Exec(transaction.ID,transaction.CreditCardID,transaction.Amount,transaction.Status,transaction.Description,transaction.Store,transaction.CreatedAt)
	if err != nil {
		return err
	}

	if transaction.Status == "approved" {
		err = repository.updateBalance(creditCard)
		if err != nil {
			return err
		}
	}
	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (repository *TransactionRepository) updateBalance(creditCard entity.CreditCard) error {
	_, err := repository.DB.Exec("update credit_cards set balance=$1 where id=$2", creditCard.Balance,creditCard.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepository) CreateCreditCard(creditCard entity.CreditCard) error {
	stmt, err := t.DB.Prepare(`insert into credit_cards(id, name, number, expiration_month,expiration_year, CVV,balance, balance_limit) 
								values($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (repository *TransactionRepository) GetCreditCard(creditCard entity.CreditCard) (entity.CreditCard, error) {
	var c entity.CreditCard
	stmt, err := repository.DB.Prepare("select id, balance, balance_limit from credit_cards where number=$1")
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}