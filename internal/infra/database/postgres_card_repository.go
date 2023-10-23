package database

import (
	"context"
	"database/sql"
	"errors"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/application/repository"
)

type PostgresCardRepository struct {
	db *sql.DB
}

func NewPostgresCardRepository(db *sql.DB) *PostgresCardRepository {
	return &PostgresCardRepository{
		db: db,
	}
}

func (r *PostgresCardRepository) Find(ctx context.Context, cardToken string) (*entity.Card, error) {
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
			return nil, app_errors.NewNotFound(repository.ErrorCardTokenIsInvalid)
		}

		return nil, app_errors.NewInternal(err)
	}

	return &card, nil
}
