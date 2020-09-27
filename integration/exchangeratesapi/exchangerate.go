package exchangeratesapi

import (
	"net/http"
	"net/url"

	"github.com/shopspring/decimal"
)

type ExchangeRateRequest struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

type RawExchangeRateResponse struct {
	Rates map[string]decimal.Decimal `json:"rates"`
	Base  string                     `json:"base"`
	Date  string                     `json:"date"`
}

func (c *Client) GetExchangeRate(req *ExchangeRateRequest) (*RawExchangeRateResponse, error) {
	var resp *RawExchangeRateResponse

	rel := &url.URL{
		Path:     "/latest",
		RawQuery: buildQuery(req),
	}

	if _, err := c.Do(
		http.MethodGet,
		rel,
		&resp,
		nil,
	); err != nil {
		return nil, err
	}

	return resp, nil
}

func buildQuery(req *ExchangeRateRequest) string {
	v := url.Values{}
	v.Set("base", req.BaseCurrency)
	v.Set("symbols", req.QuoteCurrency)

	return v.Encode()
}
