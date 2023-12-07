package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	core_errors "github.com/sesaquecruz/go-payment-processor/internal/core/errors"
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
		slog.Error(err.Error())
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
			return nil, core_errors.NewNotFoundError("card token is invalid")
		}

		slog.Error(err.Error())
		return nil, core_errors.NewInternalError(err)
	}

	return &card, nil
}
