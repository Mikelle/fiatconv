package fiatconv

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/fiatconv/integration/exchangeratesapi"
)

type FiatconvService interface {
	ConvertFromBaseToQuote(req *ExchangeRateRequest) (*decimal.Decimal, error)
}

type Client interface {
	GetExchangeRate(req *exchangeratesapi.ExchangeRateRequest) (*ExchangeRateResponse, error)
}

type Service struct {
	client Client
}

func NewService(client Client) Service {
	return Service{
		client: client,
	}
}

func (s *Service) ConvertFromBaseToQuote(req *ExchangeRateRequest) (*decimal.Decimal, *decimal.Decimal, error) {
	er := &exchangeratesapi.ExchangeRateRequest{
		BaseCurrency:  req.BaseCurrency,
		QuoteCurrency: req.QuoteCurrency,
	}

	resp, err := s.client.GetExchangeRate(er)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get exchange rate for currencies: %v", err)
	}

	res := req.Amount.Mul(resp.Rate)
	return &res, &resp.Rate, nil
}
