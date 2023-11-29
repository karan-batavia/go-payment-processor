package main

import (
	"log"

	"github.com/sesaquecruz/go-payment-processor/config"
	"github.com/sesaquecruz/go-payment-processor/di"
	"github.com/sesaquecruz/go-payment-processor/internal/acquirer"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/connection"
	"github.com/sesaquecruz/go-payment-processor/internal/infra/service"
)

func main() {
	cfg := config.GetConfig()

	db, err := connection.DBConnection(cfg.DbDsn)
	if err != nil {
		log.Fatal(err)
	}

	app := di.NewApp(
		db,
		service.PaymentWithAcquirer(acquirer.NewCielo(cfg.CieloUrl, cfg.CieloKey)),
		service.PaymentWithAcquirer(acquirer.NewRede(cfg.RedeUrl, cfg.RedeKey)),
		service.PaymentWithAcquirer(acquirer.NewStone(cfg.StoneUrl, cfg.StoneKey)),
	)

	app.Listen(":8080")
}
