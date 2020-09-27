package exchangeratesapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL *url.URL
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: mustParseEnvironmentURL(baseURL),
		httpClient: http.DefaultClient,
	}
}

func mustParseEnvironmentURL(environmentURL string) *url.URL {
	u, err := url.Parse(environmentURL)
	if err != nil {
		panic(err)
	}

	return u
}

func (c *Client) Do(method string, rel *url.URL, v interface{}, reqBody interface{}) (*http.Response, error) {
	baseURL := &(*c.baseURL)
	baseURL = baseURL.ResolveReference(rel)

	r := &http.Request{
		Method:     method,
		URL:        baseURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       baseURL.Host,
	}

	WithJSON()(r)

	switch method {
	case "POST", "PUT", "PATCH":
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	default:
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err = CheckResponse(resp); err != nil {
		return nil, err
	}

	if v == nil {
		return resp, nil
	}

	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return resp, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(bytes.NewBuffer(body)).Decode(v); err != nil && err != io.EOF {
		return nil, err
	}

	return resp, nil
}

func CheckResponse(r *http.Response) error {
	var (
		resp *ErrorResponse
		body []byte
		err  error
	)
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = r.Body.Close(); err != nil {
		return err
	}

	if len(body) == 0 {
		return &ErrorResponse{
			Message:    "response body is empty",
			StatusCode: http.StatusUnauthorized,
		}
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if r.StatusCode == http.StatusOK {
		if err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&resp); err != nil {
			return nil
		}

		if resp.Message == "" {
			return nil
		}

		return resp
	}

	return &ErrorResponse{
		Message:    string(body),
		StatusCode: r.StatusCode,
	}
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

type OptionFunc func(*http.Request)

func WithJSON() OptionFunc {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json")
	}
}