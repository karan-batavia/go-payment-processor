package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/connection"
	"github.com/sesaquecruz/go-payment-processor/test/testcontainers"

	"github.com/stretchr/testify/suite"
)

type CardRepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	db             *sql.DB
	pgContainer    *testcontainers.PostgresContainer
	cardRepository *CardRepository
}

func (s *CardRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	migrationsPath := "../../../migrations"

	pgContainer, err := testcontainers.NewPostgresContainer(ctx, migrationsPath)
	s.Require().Nil(err)

	db, err := connection.DBConnection(pgContainer.DSN)
	s.Require().Nil(err)

	s.ctx = ctx
	s.db = db
	s.pgContainer = pgContainer
	s.cardRepository = NewCardRepository(db)
}

func (s *CardRepositoryTestSuite) TestFindCards() {
	err := s.pgContainer.ClearDB()
	s.Require().Nil(err)

	err = saveTestCardData(s.db)
	s.Require().Nil(err)

	type Input struct {
		CardToken string
	}

	type Expected struct {
		Card *entity.Card
		Err  *errors.NotFoundError
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
				errors.NewNotFoundError("card token is invalid"),
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.Test, func(t *testing.T) {
			card, err := s.cardRepository.FindCard(s.ctx, tc.Input.CardToken)
			if tc.Expected.Err == nil && err == nil {
				s.Equal(tc.Expected.Card, card)
				return
			}

			var e *errors.NotFoundError
			s.Require().ErrorAs(err, &e)
			s.Equal(tc.Expected.Err.Message, e.Message)
		})
	}
}

func (s *CardRepositoryTestSuite) TearDownSuite() {
	err := s.pgContainer.TerminateContainer()
	s.Require().Nil(err)
}

func TestCardRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CardRepositoryTestSuite))
}

func saveTestCardData(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO cards (token, holder, expiration, brand)
		VALUES
			('461c9432d4d7eca7ba32b783aa22ca5c89e4f396288de5128b73b461c42d4f40', 'Bruce Wayne', '01/2025', 'VISA'),
			('7d2cd4f89ffe5374013d68c64ec104182366f786a377da1d3103db201149d3b5', 'Tony Stark', '02/2026', 'VISA'),
			('4939de8e7acf6011a9b4aa4abdd6496cec40240e418a7892723ef16c4cbb44f2', 'Peter Park', '03/2027', 'MASTERCARD'),
			('d840c6fb8401c4bbefdc4ceddc1a88f1636734bdde88c344b8969d0cd5cfdaed', 'Diana Prince', '04/2028', 'MASTERCARD'),
			('f8a8a91d9626b66a74ff7c11b5f1c2cc59a103f5b2a3119b4729a70b25304074', 'Frank Castle', '05/2029', 'AMERICAN EXPRESS'),
			('cc89fefc83d423395b11998646cc7eb7c32c04ece114d1373c3a519fbb612724', 'Natasha Romanova', '06/2030', 'AMERICAN EXPRESS');
	`)

	return err
}
