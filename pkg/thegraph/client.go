package thegraph

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	DEFAULT_URL       = "https://thegraph.bellecour.iex.ec/subgraphs/name/bellecour"
	defaulHttpTimeout = 10 * time.Second
)

var errOnTheGraph = errors.New("error while trying to fetch data from TheGraph")

// Client represents an TheGraph API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewDefaultClient creates a new TheGraph client with default URL
func NewDefaultClient() *Client {
	return NewClient(DEFAULT_URL)
}

// NewClient creates a new TheGraph client with specific URL
func NewClient(url string) *Client {
	return &Client{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: defaulHttpTimeout,
		},
	}
}

// GetVouchers fetches vouchers data from TheGraph
func (c *Client) GetVouchers() (VoucherResponse, error) {
	query := `
	{
		vouchers(orderBy: expiration, orderDirection: desc, first: 500, where: {balance_gt: 0}) {
			voucherType {
				id
				description
			}
			id
			owner {
				id
			}
			expiration
			value
			balance
		}
	}`

	var vouchers VoucherResponse
	err := c.fetchGraphQLData("/iexec-voucher", query, &vouchers)

	return vouchers, err
}

// fetchGraphQLData is a helper function to execute GraphQL queries
func (c *Client) fetchGraphQLData(endpoint, query string, result interface{}) error {
	jsonPayload, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return errOnTheGraph
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errOnTheGraph
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errOnTheGraph
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errOnTheGraph
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return errOnTheGraph
	}

	return nil
}
