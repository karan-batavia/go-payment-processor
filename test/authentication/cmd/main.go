package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/sesaquecruz/go-payment-processor/test/authentication"
)

func main() {
	fmt.Printf("\n----- public key in base64-----\n")
	fmt.Println(encodeAuthPublicKey(&authentication.PublicKey))
	fmt.Printf("------------------------------\n\n")

	app := authentication.App()
	app.Listen(":6062")
}

func encodeAuthPublicKey(publicKey *rsa.PublicKey) string {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	return publicKeyBase64
}
