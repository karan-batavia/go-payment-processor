package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

const (
	errorCardTokenIsInvalid = core_errors.Error("card token is invalid")
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (r *CardRepository) FindCard(ctx context.Context, cardToken string) (*entity.Card, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT token, holder, expiration, brand FROM cards WHERE token = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var card entity.Card
	err = stmt.QueryRowContext(ctx, cardToken).Scan(
		&card.Token,
		&card.Holder,
		&card.Expiration,
		&card.Brand,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core_errors.NewNotFoundError(errorCardTokenIsInvalid)
		}

		return nil, core_errors.NewInternalError(err)
	}

	return &card, nil
}
