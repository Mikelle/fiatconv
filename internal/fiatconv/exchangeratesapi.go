package fiatconv

import (
	"github.com/fiatconv/integration/exchangeratesapi"
)

type ExchangeRatesAPIClient struct {
	client *exchangeratesapi.Client
}

func NewExchangeRatesAPIClient(client *exchangeratesapi.Client) *ExchangeRatesAPIClient {
	return &ExchangeRatesAPIClient{
		client: client,
	}
}

func (erc *ExchangeRatesAPIClient) GetExchangeRate(req *exchangeratesapi.ExchangeRateRequest) (*ExchangeRateResponse, error) {
	resp, err := erc.client.GetExchangeRate(req)
	if err != nil {
		return nil, err
	}

	rate := &ExchangeRateResponse{Rate: resp.Rates[req.QuoteCurrency]}
	return rate, nil
}
