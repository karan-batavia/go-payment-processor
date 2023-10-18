package database

import (
	"context"
	"database/sql"
	"testing"

	app_errors "github.com/sesaquecruz/go-payment-processor/internal/application/errors"

	"github.com/sesaquecruz/go-payment-processor/internal/application/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/application/repository"
	"github.com/sesaquecruz/go-payment-processor/testcontainers"

	"github.com/stretchr/testify/suite"
)

type PostgresCardRepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	pgContainer    *testcontainers.PostgresContainer
	cardRepository *PostgresCardRepository
}

func (s *PostgresCardRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	migrationsPath := "file://../../../migrations"

	pgContainer, err := testcontainers.NewPostgresContainer(ctx, migrationsPath)
	if err != nil {
		s.FailNow(err.Error())
	}

	err = pgContainer.ClearDB()
	if err != nil {
		s.FailNow(err.Error())
	}

	db, err := sql.Open("postgres", pgContainer.DSN)
	if err != nil {
		s.FailNow(err.Error())
	}

	cardRepository := NewPostgresCardRepository(db)

	s.ctx = ctx
	s.pgContainer = pgContainer
	s.cardRepository = cardRepository
}

func (s *PostgresCardRepositoryTestSuite) TearDownSuite() {
	if err := s.pgContainer.TerminateContainer(); err != nil {
		s.FailNow(err.Error())
	}
}

func (s *PostgresCardRepositoryTestSuite) TestFindCards() {
	type Input struct {
		CardToken string
	}

	type Expected struct {
		Card *entity.Card
		Err  error
	}

	testCases := []struct {
		Test     string
		Input    *Input
		Expected *Expected
	}{
		{
			Test: "existent card",
			Input: &Input{
				"461c9432d4d7eca7ba32b783aa22ca5c89e4f396288de5128b73b461c42d4f40",
			},
			Expected: &Expected{
				entity.NewCard(
					"461c9432d4d7eca7ba32b783aa22ca5c89e4f396288de5128b73b461c42d4f40",
					"Bruce Wayne",
					"01/2025",
					"VISA",
				),
				nil,
			},
		},
		{
			Test: "existent card",
			Input: &Input{
				"4939de8e7acf6011a9b4aa4abdd6496cec40240e418a7892723ef16c4cbb44f2",
			},
			Expected: &Expected{
				entity.NewCard(
					"4939de8e7acf6011a9b4aa4abdd6496cec40240e418a7892723ef16c4cbb44f2",
					"Peter Park",
					"03/2027",
					"MASTERCARD",
				),
				nil,
			},
		},
		{
			Test: "existent card",
			Input: &Input{
				"cc89fefc83d423395b11998646cc7eb7c32c04ece114d1373c3a519fbb612724",
			},
			Expected: &Expected{
				entity.NewCard(
					"cc89fefc83d423395b11998646cc7eb7c32c04ece114d1373c3a519fbb612724",
					"Natasha Romanova",
					"06/2030",
					"AMERICAN EXPRESS",
				),
				nil,
			},
		},
		{
			Test: "non-existent card",
			Input: &Input{
				"dfadfioudafj8123jfdsfsoidfuoisdfu123zdfi",
			},
			Expected: &Expected{
				nil,
				repository.ErrorCardTokenIsInvalid,
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.Test, func(t *testing.T) {
			card, err := s.cardRepository.Find(s.ctx, tc.Input.CardToken)
			if err != nil {
				var e *app_errors.NotFound
				s.ErrorAs(err, &e)
			}

			s.ErrorIs(err, tc.Expected.Err)
			s.Equal(tc.Expected.Card, card)
		})
	}
}

func TestPostgresCardRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresCardRepositoryTestSuite))
}
