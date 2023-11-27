package main

import (
	"github.com/sesaquecruz/go-payment-processor/test/acquirer"
)

func main() {
	app := acquirer.App()
	app.Listen(":6061")
}
