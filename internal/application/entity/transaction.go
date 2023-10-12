package entity

import "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

const (
	TransactionCardIsRequired     = errors.Validation("transaction card is required")
	TransactionPurchaseIsRequired = errors.Validation("transaction purchase is required")
	TransactionStoreIsRequired    = errors.Validation("transaction store is required")
	TransactionAcquirerIsRequired = errors.Validation("transaction acquirer is required")
)

type Transaction struct {
	Card     *Card
	Purchase *Purchase
	Store    *Store
	Acquirer string
}

func NewTransaction(card *Card, purchase *Purchase, store *Store, acquirer string) *Transaction {
	return &Transaction{
		Card:     card,
		Purchase: purchase,
		Store:    store,
		Acquirer: acquirer,
	}
}

func (t *Transaction) Validate() error {
	if t.Card == nil {
		return TransactionCardIsRequired
	}

	if err := t.Card.Validate(); err != nil {
		return err
	}

	if t.Purchase == nil {
		return TransactionPurchaseIsRequired
	}

	if err := t.Purchase.Validate(); err != nil {
		return err
	}

	if t.Store == nil {
		return TransactionStoreIsRequired
	}

	if err := t.Store.Validate(); err != nil {
		return err
	}

	if t.Acquirer == "" {
		return TransactionAcquirerIsRequired
	}

	return nil
}
