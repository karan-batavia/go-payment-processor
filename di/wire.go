//go:build wireinject
// +build wireinject

package di

import (
	"crypto/rsa"
	"database/sql"

	irepository "github.com/sesaquecruz/go-payment-processor/internal/core/repository"
	iservice "github.com/sesaquecruz/go-payment-processor/internal/core/service"
	"github.com/sesaquecruz/go-payment-processor/internal/core/usecase"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/repository"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/service"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/web/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var setCardRepository = wire.NewSet(
	repository.NewCardRepository,
	wire.Bind(new(irepository.ICardRepository), new(*repository.CardRepository)),
)

var setPaymentService = wire.NewSet(
	service.NewPaymentService,
	wire.Bind(new(iservice.IPaymentService), new(*service.PaymentService)),
)

var setProcessPaymentUsecase = wire.NewSet(
	usecase.NewProcessPayment,
	wire.Bind(new(usecase.IProcessPayment), new(*usecase.ProcessPayment)),
)

var setPaymentHandler = wire.NewSet(
	handler.NewPaymentHandler,
	wire.Bind(new(handler.IPaymentHandler), new(*handler.PaymentHandler)),
)

func NewApp(db *sql.DB, authPublicKey *rsa.PublicKey, options ...service.PaymentOption) *fiber.App {
	wire.Build(
		setCardRepository,
		setPaymentService,
		setProcessPaymentUsecase,
		setPaymentHandler,
		web.InitApp,
	)

	return &fiber.App{}
}
