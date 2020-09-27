package fiatconv

import (
	"strings"

	"github.com/shopspring/decimal"
)

type FiatconvTransport interface {
	ConvertFromBaseToQuote(base, quote, amount string)
}

type Logger interface {
	Log(values ...interface{}) error
}

type Transport struct {
	logger  Logger
	service Service
}

func NewTransport(logger Logger, service Service) *Transport {
	return &Transport{
		logger: logger,
		service: service,
	}
}

func (t *Transport) ConvertFromBaseToQuote(base, quote, amount string) {
	if len(base) != 3 {
		_ = t.logger.Log("msg", "length of base is not suitable for currency", "base", base, "length", len(base))
		return
	}

	if len(quote) != 3 {
		_ = t.logger.Log("msg", "length of quote is not suitable for currency", "quote", quote, "length", len(quote))
		return
	}

	n, err := decimal.NewFromString(amount)
	if err != nil {
		_ = t.logger.Log("msg", "amount is not decimal format", "err", err, "amount", amount)
		return
	}

	req := &ExchangeRateRequest{
		BaseCurrency: strings.ToUpper(base),
		QuoteCurrency: strings.ToUpper(quote),
		Amount: n,
	}

	res, rate, err := t.service.ConvertFromBaseToQuote(req)
	if err != nil {
		_ = t.logger.Log("msg", "failed to exchange", "err", err, "base",
			req.BaseCurrency, "quote", req.QuoteCurrency, "amount", req.Amount)
		return
	}

	_ = t.logger.Log("msg", "exchange result",
		"base", req.BaseCurrency, "quote", req.QuoteCurrency,
		"exchange_rate", rate, "amount", req.Amount, "converted_amount", res)
}
