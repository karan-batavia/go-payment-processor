package acquirer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sesaquecruz/go-payment-processor/internal/core/entity"
	"github.com/sesaquecruz/go-payment-processor/internal/core/errors"
)

type Rede struct {
	name string
	url  string
	key  string
}

func NewRede(url string, key string) *Rede {
	return &Rede{
		name: "rede",
		url:  url,
		key:  key,
	}
}

func (a *Rede) Name() string {
	return a.name
}

func (a *Rede) RequestBuilder(ctx context.Context, transaction *entity.Transaction) (*http.Request, error) {
	type RedeRequest struct {
		CardToken            string   `json:"card_token"`
		CardHolder           string   `json:"card_holder"`
		CardExpiration       string   `json:"card_expiration"`
		CardBrand            string   `json:"card_brand"`
		PurchaseValue        float64  `json:"purchase_value"`
		PurchaseItems        []string `json:"purchase_items"`
		PurchaseInstallments int      `json:"purchase_installments"`
		StoreIdentification  string   `json:"store_identification"`
		StoreAddress         string   `json:"store_address"`
		StoreCep             string   `json:"store_cep"`
		StoreName            string   `json:"store_name"`
	}

	data := RedeRequest{
		CardToken:            transaction.Card.Token,
		CardHolder:           transaction.Card.Holder,
		CardExpiration:       transaction.Card.Expiration,
		CardBrand:            transaction.Card.Brand,
		PurchaseValue:        transaction.Purchase.Value,
		PurchaseItems:        transaction.Purchase.Items,
		PurchaseInstallments: transaction.Purchase.Installments,
		StoreIdentification:  transaction.Store.Identification,
		StoreAddress:         transaction.Store.Address,
		StoreCep:             transaction.Store.Cep,
		StoreName:            transaction.Acquirer.Name,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	request, err := http.NewRequest(http.MethodPost, a.url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	request.Header.Set("Api-Key", a.key)
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func (a *Rede) ResponseExtractor(response *http.Response) (*entity.Payment, error) {
	type RedeResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	var data RedeResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.NewAcquirerError(data.Code, data.Message)
	}

	payment := entity.NewPayment(data.Message)
	return payment, nil
}
