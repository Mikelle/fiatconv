package fiatconv

import "github.com/shopspring/decimal"

type ExchangeRateResponse struct {
	Rate decimal.Decimal
}

type ExchangeRateRequest struct {
	BaseCurrency  string
	QuoteCurrency string
	Amount        decimal.Decimal
}
